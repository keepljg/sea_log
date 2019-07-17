package myLoadGen

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"sea_log/logs"
	"sea_log/myLoadGen/lib"
	"sync"
	"wetalk/lib/logger"

	"math"
	"sync/atomic"
	"time"
)


type MyGenerator struct {
	caller      lib.Caller           // 调用器。
	timeoutNS   time.Duration        // 处理超时时间，单位：纳秒。
	lps         uint32               // 每秒载荷量。
	durationNS  time.Duration        // 负载持续时间，单位：纳秒。
	concurrency uint32               // 载荷并发量 即间隔时间。
	tickets     lib.GoTickets        // Goroutine票池。
	ctx         context.Context      // 上下文。
	cancelFunc  context.CancelFunc   // 取消函数。
	callCount   int64                // 调用计数。
	status      uint32               // 状态。
	resultCh    chan *lib.CallResult // 调用结果通道。
	Looping     bool
	reLoop      chan struct{}
	cond        *sync.Cond
}

// NewGenerator 会新建一个载荷发生器。
func NewGenerator(pset ParamSet) (lib.Generator, error) {

	logs.INFO("New a load generator...")
	//if err := pset.Check(); err != nil {
	//	return nil, err
	//}
	gen := &MyGenerator{
		caller:     pset.Caller,
		timeoutNS:  pset.TimeoutNS,
		lps:        pset.LPS,
		durationNS: pset.DurationNS,
		status:     lib.STATUS_ORIGINAL,
		resultCh:   pset.ResultCh,
	}
	if err := gen.init(); err != nil {
		return nil, err
	}
	return gen, nil
}

// 初始化载荷发生器。
func (gen *MyGenerator) init() error {
	var buf bytes.Buffer
	buf.WriteString("Initializing the load generator...")

	if gen.caller != nil { // 载荷测试
		// 载荷的并发量 ≈ 载荷的响应超时时间 / 载荷的发送间隔时间
		var total64 = int64(gen.timeoutNS)/int64(1e9/gen.lps) + 1
		if total64 > math.MaxInt32 {
			total64 = math.MaxInt32
		}
		gen.concurrency = uint32(total64)
		tickets, err := lib.NewGoTickets(gen.concurrency)
		if err != nil {
			return err
		}
		gen.tickets = tickets
	} else {   // 接口监控
		gen.cond = sync.NewCond(new(sync.Mutex))
		gen.Looping = true
		gen.reLoop = make(chan struct{})
	}

	buf.WriteString(fmt.Sprintf("Done. (concurrency=%d)", gen.concurrency))
	logs.INFO(buf.String())
	return nil
}

// callOne 会向载荷承受方发起一次调用。
func (gen *MyGenerator) callOne(rawReq *lib.RawReq) *lib.RawResp {
	atomic.AddInt64(&gen.callCount, 1)
	if rawReq == nil {
		return &lib.RawResp{Err: errors.New("Invalid raw request.")}
	}
	start := time.Now().UnixNano()
	resp, err := gen.caller.Call(rawReq.Req, gen.timeoutNS)
	end := time.Now().UnixNano()
	elapsedTime := time.Duration(end - start)
	var rawResp lib.RawResp
	if err != nil {
		errMsg := fmt.Sprintf("Sync Call Error: %s.", err)
		rawResp = lib.RawResp{
			Err:    errors.New(errMsg),
			Elapse: elapsedTime}
	} else {
		rawResp = lib.RawResp{
			Resp:   resp,
			Elapse: elapsedTime}
	}
	return &rawResp
}

// asyncSend 会异步地调用承受方接口。
func (gen *MyGenerator) asyncCall() {
	gen.tickets.Take()
	go func() {
		defer func() {
			if p := recover(); p != nil {
				err, ok := interface{}(p).(error)
				var errMsg string
				if ok {
					errMsg = fmt.Sprintf("Async Call Panic! (error: %s)", err)
				} else {
					errMsg = fmt.Sprintf("Async Call Panic! (clue: %#v)", p)
				}
				logs.ERROR(errMsg)
				result := &lib.CallResult{
					FuncName: gen.caller.FuncName(),
					Code:     lib.RET_CODE_FATAL_CALL,
					Msg:      errMsg}
				gen.SendResult(result)
			}
			gen.tickets.Return()
		}()
		rawReq := gen.caller.BuildReq()
		// 调用状态：0-未调用或调用中；1-调用完成；2-调用超时。
		var callStatus uint32
		timer := time.AfterFunc(gen.timeoutNS, func() {
			if !atomic.CompareAndSwapUint32(&callStatus, 0, 2) {
				return
			}
			result := &lib.CallResult{
				FuncName: gen.caller.FuncName(),
				//Req:    rawReq,
				Code:   lib.RET_CODE_WARNING_CALL_TIMEOUT,
				Msg:    fmt.Sprintf("Timeout! (expected: < %v)", gen.timeoutNS),
				Elapse: gen.timeoutNS,
			}
			gen.SendResult(result)
		})
		rawResp := gen.callOne(&rawReq)
		if !atomic.CompareAndSwapUint32(&callStatus, 0, 1) {
			return
		}
		timer.Stop()
		var result = &lib.CallResult{
			FuncName: gen.caller.FuncName(),
			Elapse:   rawResp.Elapse,
		}
		if rawResp.Err != nil {
			result.Code = lib.RET_CODE_ERROR_CALL
			result.Msg = rawResp.Err.Error()
		} else {
			result.Code = lib.RET_CODE_SUCCESS
		}
		gen.SendResult(result)
	}()
}

// sendResult 用于发送调用结果。
func (gen *MyGenerator) SendResult(result *lib.CallResult) bool {
	if atomic.LoadUint32(&gen.status) != lib.STATUS_STARTED {
		gen.printIgnoredResult(result, "stopped load generator")
		return false
	}
	select {
	case gen.resultCh <- result:
		return true
	default:
		gen.printIgnoredResult(result, "full result channel")
		return false
	}
}

// printIgnoredResult 打印被忽略的结果。
func (gen *MyGenerator) printIgnoredResult(result *lib.CallResult, cause string) {
	resultMsg := fmt.Sprintf(
		"FuncName=%s, Code=%d, Msg=%s, Elapse=%v",
		result.FuncName, result.Code, result.Msg, result.Elapse)
	logs.WARNING(fmt.Sprintf("Ignored result: %s. (cause: %s)\n", resultMsg, cause))
}

// prepareStop 用于为停止载荷发生器做准备。
func (gen *MyGenerator) prepareToStop(ctxError error) {
	logs.INFO(fmt.Sprintf("Prepare to stop load generator (cause: %s)...", ctxError))
	atomic.CompareAndSwapUint32(
		&gen.status, lib.STATUS_STARTED, lib.STATUS_STOPPING)
	logs.INFO(fmt.Sprintf("Closing result channel..."))
	close(gen.resultCh)
	atomic.StoreUint32(&gen.status, lib.STATUS_STOPPED)
}

// genLoad 会产生载荷并向承受方发送。
func (gen *MyGenerator) genLoad(throttle <-chan time.Time) {
	for {
		select {
		case <-gen.ctx.Done():
			gen.prepareToStop(gen.ctx.Err())
			return
		default:
		}
		gen.asyncCall()
		if gen.lps > 0 {
			select {
			case <-throttle:
			case <-gen.ctx.Done():
				gen.prepareToStop(gen.ctx.Err())
				return
			}
		}
	}
}

// genLoad 会产生载荷并向承受方发送。
func (gen *MyGenerator) genLoop() {
	t := time.NewTimer(gen.durationNS)
	for {
		select {
		case <-t.C:
			atomic.CompareAndSwapUint32(
				&gen.status, lib.STATUS_STARTED, lib.STATUS_STOPPING) // 暂时关闭gen
			logs.INFO("Temporary Close Gen, Prepare start Calculation")
			close(gen.resultCh)
			gen.cond.L.Lock()
			gen.cond.Wait()
			gen.cond.L.Unlock()
			gen.resultCh = make(chan *lib.CallResult, 50)
			go gen.loopCountResult(gen.cond)
			atomic.CompareAndSwapUint32(
				&gen.status, lib.STATUS_STOPPING, lib.STATUS_STARTED) // 重新启动gen
			logs.INFO("Restart Gen, Prepare Next Time Calculation")
			t.Reset(gen.durationNS)
		//case <-gen.reLoop: // 新一轮通知
		//
		//	gen.resultCh = make(chan *lib.CallResult, 50)
		//
		//	atomic.CompareAndSwapUint32(
		//		&gen.status, lib.STATUS_STOPPING, lib.STATUS_STARTED) // 重新启动gen
		//	logger.Info("Restart Gen, Prepare Next Time Calculation")
		case <-gen.ctx.Done():
			gen.prepareToStop(gen.ctx.Err())
			return
		}
	}
}

// Start 会启动载荷发生器。
func (gen *MyGenerator) Start() bool {
	logs.INFO("Starting load generator...")
	// 检查是否具备可启动的状态，顺便设置状态为正在启动
	if !atomic.CompareAndSwapUint32(
		&gen.status, lib.STATUS_ORIGINAL, lib.STATUS_STARTING) {
		if !atomic.CompareAndSwapUint32(
			&gen.status, lib.STATUS_STOPPED, lib.STATUS_STARTING) {
			return false
		}
	}
	gen.resultCh = make(chan *lib.CallResult, 50)
	if !gen.Looping {
		// 设定节流阀。
		// 初始化上下文和取消函数。
		gen.ctx, gen.cancelFunc = context.WithTimeout(
			context.Background(), gen.durationNS)
		var throttle <-chan time.Time
		if gen.lps > 0 {
			interval := time.Duration(1e9 / gen.lps)
			logs.INFO(fmt.Sprintf("Setting throttle (%v)...", interval))
			throttle = time.Tick(interval)
		}

		// 初始化调用计数。
		gen.callCount = 0

		// 设置状态为已启动。
		atomic.StoreUint32(&gen.status, lib.STATUS_STARTED)

		go func() {
			// 生成并发送载荷。
			logs.INFO("Generating loads...")
			gen.genLoad(throttle)
			logs.INFO(fmt.Sprintf("Stopped. (call count: %d)", gen.callCount))
		}()

		gen.CountResult()
		atomic.StoreUint32(&gen.status, lib.STATUS_STOPPED)
	} else {
		// 初始化上下文和取消函数。
		gen.ctx, gen.cancelFunc = context.WithCancel(context.Background())
		// 设置状态为已启动。
		atomic.StoreUint32(&gen.status, lib.STATUS_STARTED)
		go gen.genLoop()
		go gen.loopCountResult(gen.cond)

	}
	return true
}

func (gen *MyGenerator) Stop() bool {
	if !atomic.CompareAndSwapUint32(
		&gen.status, lib.STATUS_STARTED, lib.STATUS_STOPPING) {
		return false
	}
	gen.cancelFunc()
	for {
		if atomic.LoadUint32(&gen.status) == lib.STATUS_STOPPED {
			break
		}
		time.Sleep(time.Microsecond)
	}
	return true
}

func (gen *MyGenerator) Status() uint32 {
	return atomic.LoadUint32(&gen.status)
}

func (gen *MyGenerator) CallCount() int64 {
	return atomic.LoadInt64(&gen.callCount)
}

//func (gen *MyGenerator) LoopCountResult() {
//	var max, min, all float64
//	var count int
//	min = float64(gen.timeoutNS)
//	for {
//			select {
//			case r, ok := <-gen.resultCh:
//				if atomic.LoadUint32(&gen.status) == lib.STATUS_STARTED {
//					if ok {
//						e := float64(r.Elapse / time.Second)
//						all += e
//						if e > max {
//							max = e
//						}
//						if e < min {
//							min = e
//						}
//						count++
//					} else {
//						logs.INFO(fmt.Sprintf("load %d times, max spend %2f, min spend %2f, average spend %2f", count, max, min, all/float64(count)))
//						max, all, count = 0, 0, 0
//						min = float64(gen.timeoutNS)
//						atomic.CompareAndSwapUint32(
//							&gen.status, lib.STATUS_STARTED, lib.STATUS_STOPPING) // 暂时关闭gen
//						logger.Info("Temporary Close Gen, Prepare start Calculation")
//
//
//						gen.reLoop <- struct{}{}
//
//					}
//				}
//			case <-gen.ctx.Done():
//				return
//			}
//		}
//}

func (gen *MyGenerator) loopCountResult (cond *sync.Cond){
	var max, min, all float64
	var count int
	min = float64(gen.timeoutNS)
	for r := range gen.resultCh {
		e := float64(r.Elapse / time.Second)
		all += e
		if e > max {
			max = e
		}
		if e < min {
			min = e
		}
		count++
	}
	logs.INFO(fmt.Sprintf("load %d times, max spend %2.f, min spend %.2f, average spend %2.f", count, max, min, all/float64(count)))
	cond.Signal()
}



func (gen *MyGenerator) CountResult() {
	countMap := make(map[lib.RetCode]int)
	errMap := make(map[string]int)
	for r := range gen.resultCh {
		countMap[r.Code] = countMap[r.Code] + 1
		errMap[r.Msg] ++
	}

	var total int
	logger.Info("RetCode Count:")
	for k, v := range countMap {
		codePlain := lib.GetRetCodePlain(k)
		logs.INFO(fmt.Sprintf("  Code plain: %s (%d), Count: %d.\n",
			codePlain, k, v))
		total += v
	}
	for k, v := range errMap {
		logs.INFO(fmt.Sprintf("Err plain %s, count %d", k, v))
	}
	logs.INFO(fmt.Sprintf("Total: %d.\n", total))
	successCount := countMap[lib.RET_CODE_SUCCESS]
	tps := float64(successCount) / float64(gen.durationNS/1e9)
	logs.INFO(fmt.Sprintf("Loads per second: %d; Treatments per second: %f.\n", gen.lps, tps))
}

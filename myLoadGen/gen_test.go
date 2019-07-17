package myLoadGen

import (
	"fmt"
	"sea_log/myLoadGen/lib"
	"testing"
	"time"
)

// 接口荷载测试
func TestMyGenerator_Start(t *testing.T) {
	caller := lib.NewGetCaller("test", "https://www.baidu.com", nil)
	pset := ParamSet{
		Caller:     caller,
		TimeoutNS:  300 * time.Millisecond,
		LPS:        100,
		DurationNS: 10 * time.Second,
	}
	gen, err := NewGenerator(pset)
	if err != nil {
		panic(err)
	}
	gen.Start()
}

func TestCaller(t *testing.T) {
	caller := lib.NewGetCaller("test", "https://api3.feng.com/v1/flow/excellent?split=7&sort=yes", nil)
	resp, err := caller.Call([]byte{}, 1*time.Second)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(resp))
}

func TestNewGenerator(t *testing.T) {
	a := make(chan int, 10)
	go func() {
		for {
			select {
			case i, ok := <-a:
				if ok {
					fmt.Println(i)
				}
				time.Sleep(time.Second)
			}
		}
	}()
	a <- 2
	close(a)
	time.Sleep(5 * time.Second)
	a = make(chan int, 10)
	a <- 1
	a <- 2
	time.Sleep(time.Minute)
}

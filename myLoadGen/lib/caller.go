package lib

import (
	"time"
)

// Caller 表示调用器的接口。
type Caller interface {
	// 构建请求。
	BuildReq() RawReq
	// 调用。
	Call(req []byte, timeoutNS time.Duration) ([]byte, error)
	//// 检查响应。
	//CheckResp(err error) *CallResult

	FuncName() string
}

// get 调用器
type GetCaller struct {
	name string
	request *Crawler
}


func NewGetCaller(name string, url string,  userAgent map[string]string, ) Caller {
	return &GetCaller{
		name:    name,
		request: NewCrawler(userAgent, url),
	}
}
// get 请求 返回空 RawReq， post 请求 可随机返回RawReq
func (this *GetCaller) BuildReq() RawReq{
	return RawReq{}
}

func (this *GetCaller) Call(req []byte, timeoutNS time.Duration) ([]byte, error) {
	return this.request.Get(timeoutNS)
}



func (this *GetCaller) FuncName() string {
	return this.name
}



type PostCaller struct {
	name string
	request *Crawler
	body []byte
	typ string
}


// post 调用器
func NewPostCaller(name, typ string, url string, body []byte, userAgent map[string]string) Caller {
	return &PostCaller{
		name:    name,
		request: NewCrawler(userAgent, url),
		typ:typ,
		body:body,
	}
}

func (this *PostCaller) BuildReq() RawReq{
	return RawReq{
		Req: this.body,
	}
}

func (this *PostCaller) Call(req []byte, timeoutNS time.Duration) ([]byte, error) {
	if this.typ == "json" {
		return this.request.PostJson(req, timeoutNS)
	} else {
		return this.request.PostFrom(req, timeoutNS)
	}
}


func (this *PostCaller) FuncName() string {
	return this.name
}

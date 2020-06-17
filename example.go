package requests

import "fmt"

const UA = `这里填写自定义UA信息`

var Req *Request

func InitReq() {
	//初始化传入参数代表是否校验证书
	Req = Requests(false)
	Req.Headers = map[string]string{
		"User-Agent": UA,
	}
	Req.Cookies = map[string]string{

	}
	Req.Header()
	Req.Cookie()
}

func Get(urlStr string, params map[string]string) (resp *Response, err error) {
	return Req.Get(urlStr, params)
}

func Post(urlStr string, data map[string]string, options ...string) (resp *Response, err error) {
	return Req.Post(urlStr, data, options...)
}

func Put(urlStr string, data map[string]string, options ...string) (resp *Response, err error) {
	return Req.Put(urlStr, data, options...)
}

func Delete(urlStr string, data map[string]string, options ...string) (resp *Response, err error) {
	fmt.Println()
	return Req.Delete(urlStr, data, options...)
}
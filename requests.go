package requests

import (
	"encoding/json"
	"github.com/axgle/mahonia"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

//自定义UA
var ua = "Go-http-Ropon/1.2"

//请求相关
type Request struct {
	client  *http.Client
	httpReq *http.Request
	Headers map[string]string
	Cookies map[string]string
	Params  url.Values
}

//响应相关
type Response struct {
	res      *http.Response
	encoding string
	content  []byte
	text     string
}

//构造方法
func Requests() *Request {
	req := new(Request)
	req.client = &http.Client{}
	jar, _ := cookiejar.New(nil)
	req.client.Jar = jar
	return req
}

//转urlEncode
func convertUrl(data ...map[string]string) url.Values {
	urls := url.Values{}
	for _, d := range data {
		for key, value := range d {
			urls.Add(key, value)
		}
	}
	return urls
}

//请求头
func (req *Request) Header() {
	req.httpReq.Header.Add("User-Agent", ua)
	req.httpReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	for k, v := range req.Headers {
		req.httpReq.Header.Add(k, v)
	}
}

//cookie管理
func (req *Request) Cookie() {
	for k, v := range req.Cookies {
		tmp := &http.Cookie{Name: k, Value: v}
		req.httpReq.AddCookie(tmp)
	}
}

//get方法
func (req *Request) Get(urlStr string, params map[string]string) (resp *Response, err error) {
	rep, _ := http.NewRequest("GET", urlStr, nil)
	if params != nil {
		sURL, err := url.Parse(urlStr)
		if err != nil {
			return nil, err
		}
		sURL.RawQuery = convertUrl(params).Encode()
		rep.URL = sURL
	}
	req.httpReq = rep
	req.Header()
	req.Cookie()
	resp = &Response{}
	res, err := req.client.Do(rep)
	resp.res = res
	return resp, err
}

//post方法
func (req *Request) Post(urlStr string, data map[string]string, options ...string) (resp *Response, err error) {
	postData := convertUrl(data).Encode()
	//传入json数据
	if len(options) > 0 {
		postData = options[0]
	}
	rep, _ := http.NewRequest("POST", urlStr, strings.NewReader(postData))
	req.httpReq = rep
	req.Header()
	req.Cookie()
	if len(options) > 0 {
		req.httpReq.Header.Set("Content-Type", "application/json")
	}
	resp = &Response{}
	res, err := req.client.Do(rep)
	resp.res = res
	return resp, err
}

func (res *Response) Encoding(encoding string) {
	res.encoding = encoding
}

//返回文本信息
func (res *Response) Text() (text string) {
	body, _ := ioutil.ReadAll(res.res.Body)
	if res.encoding == "gbk" {
		dec := mahonia.NewDecoder("gbk")
		text = dec.ConvertString(string(body))
	} else {
		text = string(body)
	}
	res.text = text
	return
}

func (res *Response) Body() (body io.Reader) {
	if res.encoding == "gbk" {
		dec := mahonia.NewDecoder("gbk")
		body = dec.NewReader(res.res.Body)
	} else {
		body = res.res.Body
	}
	return
}

func (res *Response) Content() (content []byte) {
	var rd io.Reader
	if res.encoding == "gbk" {
		dec := mahonia.NewDecoder("gbk")
		rd = dec.NewReader(res.res.Body)
	} else {
		rd = res.res.Body
	}
	content, _ = ioutil.ReadAll(rd)
	res.content = content
	return
}

//返回json数据
func (res *Response) Json(v interface{}) error {
	if res.content == nil {
		res.Content()
	}
	return json.Unmarshal(res.content, v)
}

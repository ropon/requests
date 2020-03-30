package requests

import (
	"compress/gzip"
	"compress/zlib"
	"crypto/tls"
	"github.com/axgle/mahonia"
	jsoniter "github.com/json-iterator/go"
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
func Requests(igTls bool) *Request {
	req := new(Request)
	if igTls {
		//忽略证书校验
		req.client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}
	} else {
		req.client = &http.Client{}
	}
	jar, _ := cookiejar.New(nil)
	req.client.Jar = jar
	return req
}

//转urlencoded编码
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
	//默认urlencoded编码
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

//配置编码 默认utf8
func (res *Response) Encoding(encoding string) {
	res.encoding = encoding
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
	if res.content != nil {
		return res.content
	}
	var reader io.Reader
	//判断压缩编码
	switch res.res.Header.Get("Content-Encoding") {
	case "gzip":
		reader, _ = gzip.NewReader(res.res.Body)
	case "deflate":
		reader, _ = zlib.NewReader(res.res.Body)
	default:
		reader = res.res.Body
	}
	if res.encoding == "gbk" {
		dec := mahonia.NewDecoder("gbk")
		reader = dec.NewReader(res.res.Body)
	}
	content, _ = ioutil.ReadAll(reader)
	res.content = content
	return
}

//返回文本信息
func (res *Response) Text() (text string) {
	text = string(res.Content())
	res.text = text
	return
}

//返回jsoniter.Any 通过Get()进一步获取值
func (res *Response) Json() jsoniter.Any {
	if res.content == nil {
		res.Content()
	}
	return jsoniter.Get(res.content)
}

//响应头信息
func (res *Response) Header() map[string][]string {
	return res.res.Header
}

//响应状态码
func (res *Response) Status() int {
	return res.res.StatusCode
}

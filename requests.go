/*
Author:Ropon
Date:  2020-12-17
*/
package requests

import (
	"compress/gzip"
	"compress/zlib"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"reflect"
	"strings"
	"sync"
	"time"
)

//自定义UA
var ua = "Go-http-CodoonOps/1.2"

//请求相关
type Request struct {
	client  *http.Client
	httpReq *http.Request
	Headers map[string]string
	Cookies map[string]string
	Params  url.Values
	mutex   *sync.RWMutex
}

//响应相关
type Response struct {
	res     *http.Response
	content []byte
	text    string
}

//构造方法
func New(options ...bool) *Request {
	var igTls bool
	if len(options) > 0 {
		igTls = options[0]
	}
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
	req.httpReq = &http.Request{
		Header: make(http.Header),
	}
	req.mutex = &sync.RWMutex{}
	req.httpReq.Header.Add("User-Agent", ua)
	req.EnableCookie(true)
	req.SetTimeout(time.Second * 5)
	return req
}

//默认req
var defaultReq = New()

//转urlencoded编码
func convertUrl(data ...map[string]interface{}) url.Values {
	urls := url.Values{}
	for _, d := range data {
		for key, value := range d {
			urls.Add(key, fmt.Sprintf("%v", value))
		}
	}
	return urls
}

//设置超时时间
func (req *Request) SetTimeout(n time.Duration) {
	req.client.Timeout = n
}

func (req *Request) EnableCookie(enable bool) {
	if enable {
		jar, _ := cookiejar.New(nil)
		req.client.Jar = jar
	} else {
		req.client.Jar = nil
	}
}

//请求头
func (req *Request) Header() {
	req.mutex.Lock()
	defer req.mutex.Unlock()
	for k, v := range req.Headers {
		req.httpReq.Header.Set(k, v)
	}
}

//Cookie管理
func (req *Request) Cookie() {
	req.mutex.Lock()
	defer req.mutex.Unlock()
	for k, v := range req.Cookies {
		tmp := &http.Cookie{Name: k, Value: v}
		req.httpReq.AddCookie(tmp)
	}
}

//get方法
func (req *Request) Get(urlStr string, options ...interface{}) (resp *Response, err error) {
	rep, _ := http.NewRequest("GET", urlStr, nil)
	var paramsData string
	if len(options) > 0 {
		data := options[0]
		switch data.(type) {
		case map[string]interface{}:
			paramsData = convertUrl(data.(map[string]interface{})).Encode()
		case string:
			paramsData = data.(string)
		}
	}
	sURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	sURL.RawQuery = paramsData
	rep.URL = sURL
	rep.Header = req.httpReq.Header
	req.httpReq = rep
	resp = &Response{}
	res, err := req.client.Do(rep)
	resp.res = res
	return resp, err
}

func (req *Request) BaseReq(Method, urlStr string, options ...interface{}) (resp *Response, err error) {
	var postData string
	if len(options) > 0 {
		data := options[0]
		switch data.(type) {
		case map[string]interface{}:
			postData = convertUrl(data.(map[string]interface{})).Encode()
			req.httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		case string:
			postData = data.(string)
			req.httpReq.Header.Set("Content-Type", "application/json")
		}
	}
	rep, _ := http.NewRequest(Method, urlStr, strings.NewReader(postData))
	rep.Header = req.httpReq.Header
	req.httpReq = rep
	resp = &Response{}
	res, err := req.client.Do(rep)
	resp.res = res
	return resp, err
}

//post方法
func (req *Request) Post(urlStr string, options ...interface{}) (resp *Response, err error) {
	return req.BaseReq("POST", urlStr, options...)
}

//put方法
func (req *Request) Put(urlStr string, options ...interface{}) (resp *Response, err error) {
	return req.BaseReq("PUT", urlStr, options...)
}

//delete方法
func (req *Request) Delete(urlStr string, options ...interface{}) (resp *Response, err error) {
	return req.BaseReq("DELETE", urlStr, options...)
}

func (res *Response) Body() (body io.Reader) {
	defer res.res.Body.Close()
	body = res.res.Body
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

func (res *Response) RawJson(v interface{}) error {
	if res.content == nil {
		res.Content()
	}
	return json.Unmarshal(res.content, &v)
}

func (res *Response) Json() (Value,error) {
	if res.content == nil {
		res.Content()
	}
	return NewJson(res.content)
}

//响应头信息
func (res *Response) Header() map[string][]string {
	return res.res.Header
}

//响应状态码
func (res *Response) Status() int {
	return res.res.StatusCode
}

//响应cookie信息
func (res *Response) Cookie() []*http.Cookie {
	return res.res.Cookies()
}

func Get(urlStr string, options ...interface{}) (resp *Response, err error) {
	return defaultReq.Get(urlStr, options...)
}

func Post(urlStr string, options ...interface{}) (resp *Response, err error) {
	return defaultReq.Post(urlStr, options...)
}

func Put(urlStr string, options ...interface{}) (resp *Response, err error) {
	return defaultReq.Put(urlStr, options...)
}

func Delete(urlStr string, options ...interface{}) (resp *Response, err error) {
	return defaultReq.Delete(urlStr, options...)
}

//结构体指针转map，传入结构体指针
func StructPtr2Map(obj interface{}, tagName string) map[string]interface{} {
	tmpVal := reflect.ValueOf(obj)
	v := tmpVal.Elem()
	t := v.Type()
	var data = make(map[string]interface{})
	for i := 0; i < v.NumField(); i++ {
		data[t.Field(i).Tag.Get(tagName)] = v.Field(i).Interface()
	}
	return data
}

//结构体转map，传入结构体
func Struct2Map(obj interface{}, tagName string) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Tag.Get(tagName)] = v.Field(i).Interface()
	}
	return data
}

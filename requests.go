package requests

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"net/url"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"
)

// 自定义UA
var ua = "Go-http-Ropon/2.1"

// Request 请求相关
type Request struct {
	client  *http.Client
	httpReq *http.Request
	Params  url.Values
	mutex   *sync.RWMutex
	Headers map[string]string
	Cookies map[string]string
	Debug   bool
	BaseUrl *url.URL
}

// Response 响应相关
type Response struct {
	res     *http.Response
	content []byte
	text    string
}

// New 构造方法
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

// 默认req
var defaultReq = New()

// 转urlencoded编码
func convertUrl(data ...map[string]interface{}) url.Values {
	urls := url.Values{}
	for _, d := range data {
		for key, value := range d {
			if v, ok := value.([]string); ok {
				for _, vv := range v {
					urls.Add(key, vv)
				}
			} else if v2, ok2 := value.([]int); ok2 {
				for _, vv2 := range v2 {
					urls.Add(key, fmt.Sprintf("%v", vv2))
				}
			} else {
				urls.Add(key, fmt.Sprintf("%v", value))
			}
		}
	}
	return urls
}

// SetTimeout 设置超时时间
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

// 设置基本认证
func (req *Request) SetBasicAuth(username, password string) {
	req.httpReq.SetBasicAuth(username, password)
}

// 设置Proxy
func (req *Request) SetProxy(proxyUrl string) {
	urlProxy, _ := url.Parse(proxyUrl)
	req.client.Transport = &http.Transport{
		Proxy:           http.ProxyURL(urlProxy),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
}

func (req *Request) SetBaseUrl(baseUrl string) error {
	baseURL, err := url.Parse(baseUrl)
	if err != nil {
		return err
	}
	req.BaseUrl = baseURL
	return nil
}

func (req *Request) SetHeader(key, val string) {
	req.httpReq.Header.Set(key, val)
}

func (req *Request) RequestDebug() {
	if !req.Debug {
		return
	}
	fmt.Println("===========Go RequestDebug ============")

	// 输出基本请求信息
	fmt.Printf("Method: %s\n", req.httpReq.Method)
	fmt.Printf("URL: %s\n", req.httpReq.URL)

	// 输出请求头
	fmt.Println("Headers:")
	for key, values := range req.httpReq.Header {
		for _, value := range values {
			fmt.Printf("%s: %s\n", key, value)
		}
	}

	// 输出请求参数
	fmt.Println("Request Parameters:")

	// 处理 URL 查询参数
	queryParams := req.httpReq.URL.Query()
	if len(queryParams) > 0 {
		fmt.Println("Query Parameters:")
		for key, values := range queryParams {
			for _, value := range values {
				fmt.Printf("%s: %s\n", key, value)
			}
		}
	}

	// 处理 body 数据
	contentType := req.httpReq.Header.Get("Content-Type")

	// 使用 httputil.DumpRequestOut 来安全地获取完整请求，包括 body
	dumpedReq, err := httputil.DumpRequestOut(req.httpReq, true)
	if err == nil {
		parts := bytes.SplitN(dumpedReq, []byte("\r\n\r\n"), 2)
		if len(parts) == 2 {
			body := parts[1]
			if strings.Contains(contentType, "application/json") {
				fmt.Println("JSON Body:")
				fmt.Println(string(body))
			} else if strings.Contains(contentType, "application/x-www-form-urlencoded") {
				fmt.Println("Form Data:")
				formData, _ := url.ParseQuery(string(body))
				for key, values := range formData {
					for _, value := range values {
						fmt.Printf("%s: %s\n", key, value)
					}
				}
			} else {
				fmt.Println("Body:")
				fmt.Println(string(body))
			}
		}
	}
	// 输出 cookies
	if len(req.client.Jar.Cookies(req.httpReq.URL)) > 0 {
		fmt.Println("\nCookies:")
		for _, cookie := range req.client.Jar.Cookies(req.httpReq.URL) {
			fmt.Println(cookie)
		}
	}

	fmt.Println("===========End RequestDebug============")
}

// Header 请求头
func (req *Request) Header() {
	req.mutex.Lock()
	defer req.mutex.Unlock()
	for k, v := range req.Headers {
		req.httpReq.Header.Set(k, v)
	}
}

// Cookie 管理
func (req *Request) Cookie() {
	req.mutex.Lock()
	defer req.mutex.Unlock()
	for k, v := range req.Cookies {
		tmp := &http.Cookie{Name: k, Value: v}
		req.httpReq.AddCookie(tmp)
	}
}

func (req *Request) Do() (*Response, error) {
	res, err := req.client.Do(req.httpReq)
	if err != nil {
		return nil, err
	}
	resp := new(Response)
	resp.res = res
	return resp, nil
}

// Get get方法
func (req *Request) Get(urlStr string, options ...interface{}) (resp *Response, err error) {
	var paramsData string
	req.mutex.Lock()
	defer req.mutex.Unlock()
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
	if req.BaseUrl != nil {
		sURL = req.BaseUrl.ResolveReference(sURL)
	}
	if paramsData != "" {
		if sURL.RawQuery != "" {
			sURL.RawQuery = fmt.Sprintf(`%s&%s`, sURL.RawQuery, paramsData)
		} else {
			sURL.RawQuery = paramsData
		}
	}
	req.httpReq.Method = "GET"
	req.httpReq.URL = sURL
	req.RequestDebug()
	return req.Do()
}

func (req *Request) BaseReq(Method, urlStr string, options ...interface{}) (resp *Response, err error) {
	var postData string
	req.mutex.Lock()
	defer req.mutex.Unlock()
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
	sURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	if req.BaseUrl != nil {
		sURL = req.BaseUrl.ResolveReference(sURL)
	}
	req.httpReq.Method = Method
	req.httpReq.URL = sURL
	req.httpReq.Body = ioutil.NopCloser(strings.NewReader(postData))
	req.httpReq.ContentLength = int64(len(postData))

	req.RequestDebug()
	res, err := req.Do()

	req.httpReq.Body = nil
	req.httpReq.GetBody = nil
	req.httpReq.ContentLength = 0
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Post post方法
func (req *Request) Post(urlStr string, options ...interface{}) (resp *Response, err error) {
	return req.BaseReq("POST", urlStr, options...)
}

// Put put方法
func (req *Request) Put(urlStr string, options ...interface{}) (resp *Response, err error) {
	return req.BaseReq("PUT", urlStr, options...)
}

// Patch put方法
func (req *Request) Patch(urlStr string, options ...interface{}) (resp *Response, err error) {
	return req.BaseReq("PATCH", urlStr, options...)
}

// Delete delete方法
func (req *Request) Delete(urlStr string, options ...interface{}) (resp *Response, err error) {
	return req.BaseReq("DELETE", urlStr, options...)
}

func (res *Response) SetRes(rawRes *http.Response) {
	res.res = rawRes
	return
}

func (res *Response) Body() (body io.ReadCloser) {
	body = res.res.Body
	return
}

func (res *Response) Content() (content []byte) {
	if res.content != nil {
		return res.content
	}
	var reader io.ReadCloser
	//判断压缩编码
	switch res.res.Header.Get("Content-Encoding") {
	case "gzip":
		reader, _ = gzip.NewReader(res.Body())
	case "deflate":
		reader, _ = zlib.NewReader(res.Body())
	default:
		reader = res.Body()
	}
	defer func(reader io.ReadCloser) {
		err := reader.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
	}(reader)
	content, _ = ioutil.ReadAll(reader)
	res.content = content
	return
}

// Text 返回文本信息
func (res *Response) Text() (text string) {
	text = string(res.Content())
	res.text = text
	return
}

func (res *Response) RawJson(v interface{}) error {
	return json.Unmarshal(res.Content(), &v)
}

func (res *Response) Json() Value {
	v, err := NewJson(res.Content())
	if err != nil {
		fmt.Println(err.Error())
	}
	return v
}

func (res *Response) SaveFile(fileName string) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(res.Content())
	f.Sync()
	return err
}

// Header 响应头信息
func (res *Response) Header() http.Header {
	return res.res.Header
}

// Status 响应状态码
func (res *Response) Status() int {
	return res.res.StatusCode
}

// Cookie 响应信息
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

func Patch(urlStr string, options ...interface{}) (resp *Response, err error) {
	return defaultReq.Patch(urlStr, options...)
}

func Delete(urlStr string, options ...interface{}) (resp *Response, err error) {
	return defaultReq.Delete(urlStr, options...)
}

// StructPtr2Map 结构体指针转map，传入结构体指针
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

// Struct2Map 结构体转map，传入结构体
func Struct2Map(obj interface{}, tagName string) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Tag.Get(tagName)] = v.Field(i).Interface()
	}
	return data
}

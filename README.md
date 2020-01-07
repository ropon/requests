#requests
```go
#初始化请求
request = requests.Requests()

#构建header
request.Headers = map[string]string{"key":val}

#构建cookie
request.Cookies = map[string]string{"key":val}

#具体请求可以单独再指定header头信息和cookie信息
request.Headers["key"] = val
request.Cookies["key"] = val

#有参数get请求
#www.ropon.top?key=val
res, err := request.Get("www.ropon.top", map[string]string{"key":val})
if err != nil {
    log.Fatal(err)
}

#无参数get请求
res, err := request.Get("www.ropon.top", nil)
if err != nil {
    log.Fatal(err)
}

#post请求 默认使用urlencoded编码
res, err := request.Post("www.ropon.top", map[string]string{"key":val})
if err != nil {
    log.Fatal(err)
}

#post 传入json请求
jsonStr := `{"key": val}`
res, err := request.Post("www.ropon.top", nil, jsonStr)
if err != nil {
    log.Fatal(err)
}

#返回对象公共参数
#默认是utf-8 若网页编码是gbk 使用
res.Encoding("gbk")
#获取文本信息
res.Text()
#获取Json
res.Json()
```

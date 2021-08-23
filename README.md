Requests
=======

下载
------------

```
go get -u github.com/Ropon/requests
```

使用
-------

```go
//构建一个新对象
req := requests.New()
req.Get("url")
//或者直接使用默认对象请求
requests.Get("url")

//构建请求头
req.Headers = map[string]string{
    "User-Agent": `这里填写自定义UA信息`,
}
req.Header()

//构建Cookie
req.Cookies = map[string]string{
    "key": "val"
}
req.Cookie()
```

**GET**:

```go
//无参数Get请求
res, err := req.Get("https://www.ropon.top")
//错误处理
if err != nil {
    log.Fatal(err)
}
fmt.Println(res.Text())

//有参数Get请求
queryData := map[string]interface{}{
    "key": "val",
    "key2": 123
}
res, err := req.Get("https://www.ropon.top", queryData)
//错误处理
if err != nil {
    log.Fatal(err)
}
fmt.Println(res.Text())
```

**POST**:

```go
//默认urlencode编码
postData := map[string]interface{}{
    "key": "val",
	"key2": 123
}
res, err := req.Post("https://www.ropon.top", postData)
//错误处理
if err != nil {
    log.Fatal(err)
}
fmt.Println(res.Text())

//请求传json
postJsonStr := `{"key": "val"}`
res, err := req.Post("https://www.ropon.top", postJsonStr)
//错误处理
if err != nil {
    log.Fatal(err)
}
fmt.Println(res.Text())
```

**PUT**:

```go
//默认urlencode编码
putData := map[string]interface{}{
    "key": "val",
	"key2": 123
}
res, err := req.Put("https://www.ropon.top", putData)
//错误处理
if err != nil {
    log.Fatal(err)
}
fmt.Println(res.Text())

//请求传json
putJsonStr := `{"key": "val"}`
res, err := req.Put("https://www.ropon.top", putJsonStr)
//错误处理
if err != nil {
    log.Fatal(err)
}
fmt.Println(res.Text())
```

**PATCH**:

```go
//默认urlencode编码
patchData := map[string]interface{}{
    "key": "val",
	"key2": 123
}
res, err := req.Patch("https://www.ropon.top", patchData)
//错误处理
if err != nil {
    log.Fatal(err)
}
fmt.Println(res.Text())

//请求传json
patchJsonStr := `{"key": "val"}`
res, err := req.Patch("https://www.ropon.top", patchJsonStr)
//错误处理
if err != nil {
    log.Fatal(err)
}
fmt.Println(res.Text())
```

**DELETE**:

```go
//默认urlencode编码
deleteData := map[string]interface{}{
    "key": "val",
    "key2": 123
}
res, err := req.Delete("https://www.ropon.top", deleteData)
//错误处理
if err != nil {
    log.Fatal(err)
}
fmt.Println(res.Text())

//请求传json
deleteJsonStr := `{"key": "val"}`
res, err := req.Delete("https://www.ropon.top", deleteJsonStr)
//错误处理
if err != nil {
    log.Fatal(err)
}
fmt.Println(res.Text())
```

**Cookies**:

```go
//单独带cookie请求
req.Cookies["key1"] = "val1"
req.Cookie()
res, err := req.Get("https://www.ropon.top")
```

**Headers**:

```go
//单独带header请求
req.Headers["key1"] = "val1"
req.Header()
res, err := req.Get("https://www.ropon.top")
```

**Res**:

```go
//获取文本信息
res.Text()
//获取Json，反馈requests.Value
res.Json()
res.Json().Get("data", "svc_list").String()
//获取响应头
res.Header()
//获取状态码
res.Status()
```

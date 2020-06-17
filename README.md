Requests
=======

下载
------------

```
go get -v github.com/Ropon/requests
```


使用
-------

```go
import (
    "github.com/Ropon/requests"
)

//参数igTls 是否忽略校验证书
Req = Requests(false)
//构建请求头
Req.Headers = map[string]string{
    "User-Agent": `这里填写自定义UA信息`,
}
Req.Header()
//构建Cookie
Req.Cookies = map[string]string{
    "key": "val"
}
Req.Cookie()
```

**GET**:

```go
//无参数Get请求
res, err := Req.Get("https://www.ropon.top", nil)
//错误处理
if err != nil {
    log.Fatal(err)
}
fmt.Println(res.Text())

//有参数Get请求
queryData := map[string]string{
    "key": "val"
}
res, err := Req.Get("https://www.ropon.top", queryData)
//错误处理
if err != nil {
    log.Fatal(err)
}
fmt.Println(res.Text())
```

**POST**:

```go
//默认urlencode编码
postData := map[string]string{
    "key": "val"
}
res, err := Req.Post("https://www.ropon.top", postData)
//错误处理
if err != nil {
    log.Fatal(err)
}
fmt.Println(res.Text())

//Post请求传json
postJsonStr := `{"key": "val"}`
res, err := Req.Post("https://www.ropon.top", nil, postJsonStr)
//错误处理
if err != nil {
    log.Fatal(err)
}
fmt.Println(res.Text())
```

**PUT**:

```go
//默认urlencode编码
putData := map[string]string{
    "key": "val"
}
res, err := Req.Put("https://www.ropon.top", putData)
//错误处理
if err != nil {
    log.Fatal(err)
}
fmt.Println(res.Text())

//Post请求传json
putJsonStr := `{"key": "val"}`
res, err := Req.Put("https://www.ropon.top", nil, putJsonStr)
//错误处理
if err != nil {
    log.Fatal(err)
}
fmt.Println(res.Text())
```

**DELETE**:

```go
//默认urlencode编码
deleteData := map[string]string{
    "key": "val"
}
res, err := Req.Delete("https://www.ropon.top", deleteData)
//错误处理
if err != nil {
    log.Fatal(err)
}
fmt.Println(res.Text())

//Post请求传json
deleteJsonStr := `{"key": "val"}`
res, err := Req.Delete("https://www.ropon.top", nil, deleteJsonStr)
//错误处理
if err != nil {
    log.Fatal(err)
}
fmt.Println(res.Text())
```

**Cookies**:

```go
//单独带cookie请求
Req.Cookies["key1"] = "val1"
res, err := Req.Get("https://www.ropon.top", nil)
```

**Headers**:

```go
//单独带header请求
Req.Headers["key1"] = "val1"
res, err := Req.Get("https://www.ropon.top", nil)
```

**Res**:

```go
//默认utf8编码，若使用gbk编码
res.Encoding("gbk")
//获取文本信息
res.Text()
//获取Json，反馈jsoniter.Any
res.Json()
//获取响应头
res.Header()
//获取状态码
res.Status()
```

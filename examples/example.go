package main

import (
	"fmt"
	"github.com/ropon/requests"
)

func main() {
	//构建一个新对象
	//req := requests.New()
	//req.Get("url")
	//或者直接使用默认对象请求
	//requests.Get("url")

	apiUrl := "https://ip.ropon.top"
	queryData := map[string]interface{}{
		"type": "json",
	}
	res, err := requests.Get(apiUrl, queryData)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//输出文本
	fmt.Println("text:", res.Text())
	//解析JSON
	code := res.Json().Get("code").String()
	//直接转数组
	zone := res.Json().Get("data", "zone").String()
	fmt.Println(code, zone)
	//循环
	userList := res.Json().Get("data")
	for i := 0; i < userList.Size(); i++ {
		fmt.Println(userList.Get(i).Get("email").String())
	}
	//post方法
	postData := map[string]interface{}{
		"type": "json",
	}
	res, err = requests.Post(apiUrl, postData)
	if err != nil {
		return
	}
	//或json
	jsonData := fmt.Sprintf(`{"name": "%s", "age": %d}`, "ropon", 18)
	res, err = requests.Post(apiUrl, jsonData)
	if err != nil {
		return
	}
	//res同上
}

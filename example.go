package requests

import (
	"fmt"
	"github.com/Ropon/requests"
)

func main() {
	//构建一个新对象
	//req := requests.New()
	//req.Get("url")
	//或者直接使用默认对象请求
	//requests.Get("url")

	apiUrl := "https://ip.ropon.top"
	data := map[string]interface{}{
		"type": "all",
	}
	res, err := requests.Get(apiUrl, data)
	if err != nil {
		return
	}
	//输出文本
	//fmt.Println(res.Text())
	//解析JSON
	status := res.Json().Get("status").String()
	//直接转数组
	nameList := res.Json().Get("data").StringArray("name")
	fmt.Println(status, nameList)
	//循环
	userList := res.Json().Get("data")
	for i := 0; i < userList.Size(); i++ {
		fmt.Println(userList.Get(i).Get("email").String())
	}
	//post方法
	data = map[string]interface{}{
		"type": "all",
	}
	//或json
	data = fmt.Sprintf(`{"name": "ropon", "age": 18}`)
	res, err = requests.Post(apiUrl, data)
	if err != nil {
		return
	}
	//res同上
}
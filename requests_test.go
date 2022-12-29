package requests

import "testing"

var req = New()

func TestGet(t *testing.T) {
	req.Debug = true
	//res, err := Get("https://httpbin.org/get?a=1&b=2","c=3&d=4")
	res, err := req.Get("https://httpbin.org/get?a=1&b=2", map[string]interface{}{
		"key1": "val1",
		"key2": 11,
		"key3": []string{"val31", "val32"},
		"key4": []int{41, 42},
	})
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(res.Text())
}

func TestPost(t *testing.T) {
	//res, err := Post("https://httpbin.org/post",`{"name":"ropon","age":18}`)
	res, err := req.Post("https://httpbin.org/post?arg1=123&arg2=456", map[string]interface{}{
		"key1": "val1",
		"key2": 22,
		"key3": []string{"val31", "val32"},
		"key4": []int{41, 42},
	})
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(res.Text())
}

func TestJson(t *testing.T) {
	//res, err := Post("https://httpbin.org/post",`{"name":"ropon","age":18}`)
	res, err := req.Post("https://httpbin.org/post?arg1=123&arg2=456", map[string]interface{}{
		"key1": "val1",
		"key2": 22,
	})
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(res.Json().Get("args"))
}

func TestProxy(t *testing.T) {
	req.SetProxy("http://127.0.0.1:7890")
	res, err := req.Get("https://www.v2ex.com")
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(res.Text())
}

func TestDownLoad(t *testing.T) {
	res, err := req.Get("https://studygolang.com/dl/golang/go1.18.3.src.tar.gz")
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(res.SaveFile("./go.tar.gz"))
}

func TestHeader(t *testing.T) {
	res, err := req.Get("https://httpbin.org/get")
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(res.Header().Get("Date"))
}

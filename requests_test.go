package requests

import "testing"

func TestGet(t *testing.T) {
	//res, err := Get("https://httpbin.org/get?a=1&b=2","c=3&d=4")
	res, err := Get("https://httpbin.org/get?a=1&b=2", map[string]interface{}{
		"key1": "val1",
		"key2": 11,
	})
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(res.Text())
}

func TestPost(t *testing.T) {
	//res, err := Post("https://httpbin.org/post",`{"name":"ropon","age":18}`)
	res, err := Post("https://httpbin.org/post?arg1=123&arg2=456", map[string]interface{}{
		"key1": "val1",
		"key2": 22,
	})
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(res.Text())
}

func TestJson(t *testing.T) {
	//res, err := Post("https://httpbin.org/post",`{"name":"ropon","age":18}`)
	res, err := Post("https://httpbin.org/post?arg1=123&arg2=456", map[string]interface{}{
		"key1": "val1",
		"key2": 22,
	})
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(res.Json().Get("args"))
}

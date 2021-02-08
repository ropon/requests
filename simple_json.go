/*
Author:Ropon
Date:  2021-01-06
*/
package requests

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Value struct {
	Data   interface{}
	Exists bool
}

func NewJson(data []byte) (Value, error) {
	v := Value{}
	err := json.Unmarshal(data, &v.Data)
	if err == nil {
		v.Exists = true
	} else {
		return v, err
	}
	return v, nil
}

func (v Value) Get(path ...interface{}) Value {
	if v.Data == nil || !v.Exists {
		return v
	}
	for _, key := range path {
		if v.Data == nil || !v.Exists {
			break
		}
		switch v.Data.(type) {
		case map[string]interface{}:
			v.Data = v.Data.(map[string]interface{})[key.(string)]
		case []interface{}:
			varray := v.Data.([]interface{})
			if key.(int) < len(varray) {
				v.Data = varray[key.(int)]
			} else {
				v.Data = nil
			}
		}
	}
	return v
}

func (v Value) String() string {
	if v.Data == nil {
		return ""
	}
	return fmt.Sprintf(`%v`, v.Data)
}

func (v Value) Int64() int64 {
	if v.Data == nil || !v.Exists {
		return 0
	}
	i, _ := strconv.ParseInt(strings.ToLower(v.String()), 10, 64)
	return i
}

func (v Value) Float64() float64 {
	if v.Data == nil || !v.Exists {
		return 0
	}
	f, _ := strconv.ParseFloat(strings.ToLower(v.String()), 64)
	return f
}

func (v Value) Bool() bool {
	if v.Data == nil || !v.Exists {
		return false
	}
	b, _ := strconv.ParseBool(strings.ToLower(v.String()))
	return b
}

func (v Value) Time(layout string, options ...string) time.Time {
	if v.Data == nil || !v.Exists {
		return time.Now()
	}
	if len(options) > 0 && options[0] == "utc" {
		tm, err := time.ParseInLocation(layout, v.String(), time.UTC)
		if err != nil {
			return time.Now()
		}
		return tm.Local()
	}
	tm, err := time.ParseInLocation(layout, v.String(), time.Local)
	if err != nil {
		return time.Now()
	}
	return tm.Local()
}

func (v Value) Size() int {
	if v.Data == nil || !v.Exists {
		return 0
	}
	switch v.Data.(type) {
	case []interface{}:
		return len(v.Data.([]interface{}))
	}
	return 0
}

func (v Value) Map() map[string]interface{} {
	tmpMap := make(map[string]interface{})
	if v.Data == nil || !v.Exists {
		return tmpMap
	}
	switch v.Data.(type) {
	case map[string]interface{}:
		tmpMap = v.Data.(map[string]interface{})
	}
	return tmpMap
}

func (v Value) StringArray(path ...interface{}) []string {
	var tmpArray []string
	if v.Data == nil || !v.Exists {
		return tmpArray
	}
	if v.Size() > 0 {
		for i := 0; i < v.Size(); i++ {
			if v.Get(i).Get(path...).String() == "" {
				break
			}
			tmpArray = append(tmpArray, v.Get(i).Get(path...).String())
		}
	}
	return tmpArray
}

package phpserialize

import "fmt"
import (
	"errors"
	"log"
	"strconv"
	"strings"
)

func Serialize() string {
	return ""
}

func Unserialize(str string) (*Serial, error) {
	s := &Serial{}
	_, _, value := unserialize(str, 0)
	fmt.Println("value:", value)
	s.data = value
	return s, nil
}

type Serial struct {
	data interface{}
}

func unserialize(data string, offset int) (string, int, interface{}) {
	dtype := strings.ToLower(data[offset : offset+1])
	dataoffset := offset + 2
	var value interface{}
	switch dtype {
	case "i":
		l, s, _ := read_until(data, dataoffset, ";")
		value = s
		dataoffset += l + 1
	case "b":
		l, s, _ := read_until(data, dataoffset, ";")
		if s == "1" {
			value = true
		} else {
			value = false
		}
		dataoffset += l + 1
	case "d":
		fmt.Println("d")
	case "n":
		value = nil
	case "s":
		l, s, _ := read_until(data, dataoffset, ":")
		dataoffset += (l + 2)
		sl, _ := strconv.Atoi(s)
		l, s = read_chrs(data, dataoffset, sl)
		dataoffset += (l + 2)
		value = s
	case "a":
		chrs, keys, _ := read_until(data, dataoffset, ":")
		dataoffset += chrs + 2
		keys_int, _ := strconv.Atoi(keys)
		t := make(map[string]interface{})
		for i := 0; i < keys_int; i++ {
			_, kchrs, key := unserialize(data, dataoffset)
			dataoffset += kchrs
			_, vchrs, va := unserialize(data, dataoffset)
			dataoffset += vchrs
			t[key.(string)] = va
		}
		value = t
	default:
		fmt.Println("d")
	}
	return dtype, dataoffset - offset, value

}

func read_until(data string, offset int, stopchr string) (int, string, error) {
	chr := data[offset]

	i := 2
	var buf []byte
	for chr != stopchr[0] {
		if (i + offset) > len(data) {
			return 0, "", errors.New("Invalid")
		}
		buf = append(buf, chr)
		chr = data[offset+(i-1)]
		i++
	}
	return len(buf), string(buf), nil
}

func read_chrs(data string, offset int, length int) (int, string) {
	// fmt.Printf("data:%s, offset:%d, length: %d\n", data, offset, length)
	var buf []byte
	for i := offset; i < offset+length; i++ {
		buf = append(buf, data[i])
	}
	return len(buf), string(buf)
}

// Get returns a pointer to a new `Json` object
// for `key` in its `map` representation
//
// useful for chaining operations (to traverse a nested JSON):
//    js.Get("top_level").Get("dict").Get("value").Int()
func (j *Serial) Get(key string) *Serial {
	m, err := j.Map()
	if err == nil {
		if val, ok := m[key]; ok {
			return &Serial{val}
		}
	}
	return &Serial{nil}
}

// GetPath searches for the item as specified by the branch
// without the need to deep dive using Get()'s.
//
//   js.GetPath("top_level", "dict")
func (j *Serial) GetPath(branch ...string) *Serial {
	jin := j
	for i := range branch {
		m, err := jin.Map()
		if err != nil {
			return &Serial{nil}
		}
		if val, ok := m[branch[i]]; ok {
			jin = &Serial{val}
		} else {
			return &Serial{nil}
		}
	}
	return jin
}

// GetIndex resturns a pointer to a new `Json` object
// for `index` in its `array` representation
//
// this is the analog to Get when accessing elements of
// a json array instead of a json object:
//    js.Get("top_level").Get("array").GetIndex(1).Get("key").Int()
func (j *Serial) GetIndex(index int) *Serial {
	a, err := j.Array()
	if err == nil {
		if len(a) > index {
			return &Serial{a[index]}
		}
	}
	return &Serial{nil}
}

// CheckGet returns a pointer to a new `Json` object and
// a `bool` identifying success or failure
//
// useful for chained operations when success is important:
//    if data, ok := js.Get("top_level").CheckGet("inner"); ok {
//        log.Println(data)
//    }
func (j *Serial) CheckGet(key string) (*Serial, bool) {
	m, err := j.Map()
	if err == nil {
		if val, ok := m[key]; ok {
			return &Serial{val}, true
		}
	}
	return nil, false
}

// Map type asserts to `map`
func (j *Serial) Map() (map[string]interface{}, error) {
	if m, ok := (j.data).(map[string]interface{}); ok {
		return m, nil
	}
	return nil, errors.New("type assertion to map[string]interface{} failed")
}

// Array type asserts to an `array`
func (j *Serial) Array() ([]interface{}, error) {
	if a, ok := (j.data).([]interface{}); ok {
		return a, nil
	}
	return nil, errors.New("type assertion to []interface{} failed")
}

// Bool type asserts to `bool`
func (j *Serial) Bool() (bool, error) {
	if s, ok := (j.data).(bool); ok {
		return s, nil
	}
	return false, errors.New("type assertion to bool failed")
}

// String type asserts to `string`
func (j *Serial) String() (string, error) {
	if s, ok := (j.data).(string); ok {
		return s, nil
	}
	return "", errors.New("type assertion to string failed")
}

// Float64 type asserts to `float64`
func (j *Serial) Float64() (float64, error) {
	if i, ok := (j.data).(float64); ok {
		return i, nil
	}
	return -1, errors.New("type assertion to float64 failed")
}

// Int type asserts to `float64` then converts to `int`
func (j *Serial) Int() (int, error) {
	if f, ok := (j.data).(float64); ok {
		return int(f), nil
	}

	return -1, errors.New("type assertion to float64 failed")
}

// Int type asserts to `float64` then converts to `int64`
func (j *Serial) Int64() (int64, error) {
	if f, ok := (j.data).(float64); ok {
		return int64(f), nil
	}

	return -1, errors.New("type assertion to float64 failed")
}

// Bytes type asserts to `[]byte`
func (j *Serial) Bytes() ([]byte, error) {
	if s, ok := (j.data).(string); ok {
		return []byte(s), nil
	}
	return nil, errors.New("type assertion to []byte failed")
}

// StringArray type asserts to an `array` of `string`
func (j *Serial) StringArray() ([]string, error) {
	arr, err := j.Array()
	if err != nil {
		return nil, err
	}
	retArr := make([]string, 0, len(arr))
	for _, a := range arr {
		s, ok := a.(string)
		if !ok {
			return nil, err
		}
		retArr = append(retArr, s)
	}
	return retArr, nil
}

// MustArray guarantees the return of a `[]interface{}` (with optional default)
//
// useful when you want to interate over array values in a succinct manner:
//		for i, v := range js.Get("results").MustArray() {
//			fmt.Println(i, v)
//		}
func (j *Serial) MustArray(args ...[]interface{}) []interface{} {
	var def []interface{}
	switch len(args) {
	case 0:
		break
	case 1:
		def = args[0]
	default:
		log.Panicf("MustArray() received too many arguments %d", len(args))
	}

	a, err := j.Array()
	if err == nil {
		return a
	}

	return def
}

// MustMap guarantees the return of a `map[string]interface{}` (with optional default)
//
// useful when you want to interate over map values in a succinct manner:
//		for k, v := range js.Get("dictionary").MustMap() {
//			fmt.Println(k, v)
//		}
func (j *Serial) MustMap(args ...map[string]interface{}) map[string]interface{} {
	var def map[string]interface{}
	switch len(args) {
	case 0:
		break
	case 1:
		def = args[0]
	default:
		log.Panicf("MustMap() received too many arguments %d", len(args))
	}

	a, err := j.Map()
	if err == nil {
		return a
	}

	return def
}

// MustString guarantees the return of a `string` (with optional default)
//
// useful when you explicitly want a `string` in a single value return context:
//     myFunc(js.Get("param1").MustString(), js.Get("optional_param").MustString("my_default"))
func (j *Serial) MustString(args ...string) string {
	var def string

	switch len(args) {
	case 0:
		break
	case 1:
		def = args[0]
	default:
		log.Panicf("MustString() received too many arguments %d", len(args))
	}

	s, err := j.String()
	if err == nil {
		return s
	}

	return def
}

// MustInt guarantees the return of an `int` (with optional default)
//
// useful when you explicitly want an `int` in a single value return context:
//     myFunc(js.Get("param1").MustInt(), js.Get("optional_param").MustInt(5150))
func (j *Serial) MustInt(args ...int) int {
	var def int

	switch len(args) {
	case 0:
		break
	case 1:
		def = args[0]
	default:
		log.Panicf("MustInt() received too many arguments %d", len(args))
	}

	i, err := j.Int()
	if err == nil {
		return i
	}

	return def
}

// MustFloat64 guarantees the return of a `float64` (with optional default)
//
// useful when you explicitly want a `float64` in a single value return context:
//     myFunc(js.Get("param1").MustFloat64(), js.Get("optional_param").MustFloat64(5.150))
func (j *Serial) MustFloat64(args ...float64) float64 {
	var def float64

	switch len(args) {
	case 0:
		break
	case 1:
		def = args[0]
	default:
		log.Panicf("MustFloat64() received too many arguments %d", len(args))
	}

	i, err := j.Float64()
	if err == nil {
		return i
	}

	return def
}

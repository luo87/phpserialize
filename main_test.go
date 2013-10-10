package phpserialize

import (
	"fmt"
	"reflect"
	"testing"
)

func TestSerialize(t *testing.T) {
	var x float64 = 3.4
	v := reflect.ValueOf(x)
	fmt.Println("type:", v.Type())
	fmt.Println("kind is float64:", v.Kind() == reflect.Float64)
	fmt.Println("value:", v.Float())
}

func TestUnserialize(t *testing.T) {
	Unserialize(`a:2:{i:232;s:10:"2323 dfsdf";s:3:"sdf";a:3:{i:0;i:2;i:1;i:3;i:2;i:4;}}`)

	data, _ := Unserialize(`a:2:{i:232;s:17:"2323 中文 dfsdf";s:3:"sdf";a:3:{i:0;i:2;i:1;i:3;i:2;i:4;}}`)
	fmt.Println(data.data)
	fmt.Println(data.Get("232").String())
}

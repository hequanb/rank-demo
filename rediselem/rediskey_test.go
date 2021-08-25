package rediselem

import (
	"fmt"
	"reflect"
	"testing"
)

func TestBuildKeys(t *testing.T) {
	BuildKeys(UserPrefix, "1", "", "3")
	var a  = []interface{}{"123","456"}
	fmt.Println(reflect.ValueOf(a).Type())
	v := reflect.ValueOf(a)
	vv := v.Index(0)
	fmt.Println(vv.Kind().String())
	fmt.Println(vv.Type().String())
}

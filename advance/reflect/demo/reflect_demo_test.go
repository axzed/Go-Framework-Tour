package demo

import (
	"fmt"
	"reflect"
	"testing"
)

func OutputFields(t testing.T) {

}

func TestReflectPanic(t *testing.T) {
	typ := reflect.TypeOf(&User{})
	if typ.Kind() == reflect.Struct {
		fmt.Println("结构体")
	} else if typ.Kind() == reflect.Ptr {
		fmt.Println("指针")
	}

}

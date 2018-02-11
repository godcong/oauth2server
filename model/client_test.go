package model

import (
	"fmt"
	"reflect"
	"testing"
)

func TestClientType(t *testing.T) {
	fmt.Println(reflect.TypeOf(CT_ADMIN))
	fmt.Println(reflect.TypeOf(CT_NONE))
	fmt.Println(reflect.TypeOf(CT_INNER))
	fmt.Println(reflect.TypeOf(CT_OUTER))

}

package configo

import (
	"fmt"
	"testing"
)

func TestNewConfig(t *testing.T) {
	t.Log(NewConfig("/", TYPE_DEFAULT) != nil)
	t.Log(NewConfig("/") != nil)
}

func TestSome(t *testing.T) {
	p := make(Property, 10)
	p["hhh"] = "hahaha"
	fmt.Println(p)
}

func TestConfig_Get(t *testing.T) {
	t.Log(Get("some"))
}

func TestConfig_Load(t *testing.T) {
	if err := config.Load(); err != nil {
		t.Log("load error", err)
	}
}
func TestGetSystemSeparator(t *testing.T) {
	t.Log(SYSTEM_SEPARATOR)
}

func ExampleGet() {
	fmt.Println(Get("example"))
}

func ExampleNewConfig() {
	NewConfig("config.eee").Properties()
}

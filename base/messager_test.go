package base

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestTranToListDesc(t *testing.T) {
	fmt.Println(strings.Compare("mobile", "servicename"))
	fmt.Println(strings.Compare("mobile", "clientid"))
	fmt.Println(strings.Compare("mobile", "timestamp"))
	fmt.Println(strings.Compare("mobile", "content"))
	fmt.Println(strings.Compare("clientid", "content"))
}

func TestTransToListByDesc(t *testing.T) {
	m := make(map[string]string)
	m["mobile"] = strconv.Itoa(13058750423)
	m["servicename"] = "mana_register"
	m["clientid"] = "clientid111"
	m["timestamp"] = time.Now().String()
	m["content"] = ""

	list := TransToListByDesc(m)
	fmt.Println(list.Len(), "!!!")
	ele := list.Front()
	for ; ele != nil; ele = ele.Next() {
		fmt.Println(ele.Value)
	}
	fmt.Println()
}

func TestSendMessage(t *testing.T) {
	m := make(map[string]string)
	m["mobile"] = strconv.Itoa(13058750423)
	m["servicename"] = "bingo_register"
	m["clientid"] = "clientid111"
	m["timestamp"] = strconv.FormatInt(time.Now().Unix(), 10)
	m["content"] = "{\"code\":\"1234\",\"product\":\"玛娜爱情\"}"
	m["type"] = "1"
	rmap := SendMessage(MT_SEND, m)
	fmt.Println(rmap)
}

func TestCheckMessage(t *testing.T) {
	m := make(map[string]string)
	m["mobile"] = strconv.Itoa(13058750423)
	//m["servicename"] = "bingo_register"
	m["clientid"] = "clientid111"
	m["timestamp"] = strconv.FormatInt(time.Now().Unix(), 10)
	//m["content"] = "{\"code\":\"1234\",\"product\":\"玛娜爱情\"}"
	m["code"] = "9999"
	//CheckMessage(m)
	rmap := SendMessage(MT_CHECK, m)
	fmt.Println(rmap)
}

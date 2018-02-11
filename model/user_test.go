package model

import "testing"

var (
	T = []string{
		"abc",
		"a123",
		"a123_",
		"1a123_",
		"_1a123_",
		"a_1_",
		"abc@",
		"abc@123.com",
		"@abc@123.com",
		"123@abc@123.com",
		"13058750425@111",
		"13058750423",
		"130587504253",
	}
)

func TestVerifyAccountType(t *testing.T) {

	t.Log("Verify UName:")
	for _, v := range T {
		t.Log(VerifyUsername(v))
	}

	t.Log("Verify Mail:")
	for _, v := range T {
		t.Log(VerifyMail(v))
	}

	t.Log("Verify Mobile:")
	for _, v := range T {
		t.Log(VerifyMobile(v))
	}

}

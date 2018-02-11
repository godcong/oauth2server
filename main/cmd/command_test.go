package main

import (
	"fmt"
	"testing"

	"github.com/godcong/oauth2server/base"
	uuid "github.com/satori/go.uuid"
)

func TestMigrate(t *testing.T) {
	Migrate()
	t.Log("TestMigrate OK")
}

func TestRandomUser(t *testing.T) {
	RandomUser("100")
	t.Log("TestRandomUser OK")
}

func TestClientMake(t *testing.T) {
	ClientMake("local")
	fmt.Println(uuid.NewV1().Bytes())
	str, e := uuid.NewV1().MarshalBinary()
	fmt.Println(string(str), e)
	str1, e := uuid.NewV1().MarshalText()
	fmt.Println(string(str1), e)
	str2, e := uuid.NewV1().MarshalText()
	fmt.Println(string(str2), e)
	str3 := base.GenerateRandomString(16, base.T_RAND_LOWERNUM)
	fmt.Println(string(str3))
}

//
//func TestTransfer(t *testing.T) {
//	model.CFlag = false
//	config := configo.NewConfig(`D:\Godcong\Workspace\g7n3\src\g7n3.com\hamster\gautu\config.env`)
//	config.Load()
//	f, err1 := config.Get("mysql2")
//
//	if err1 != nil {
//		panic(err1)
//	}
//	dfrom := connectDb(f)
//	to, err2 := config.Get("mysql")
//	if err2 != nil {
//		panic(err2)
//	}
//	dto := connectDb(to)
//	Transfer(dfrom, dto)
//
//}

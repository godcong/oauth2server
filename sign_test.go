package oauth2server

import (
	"fmt"
	"log"
	"net/url"
	"testing"

	"os"

	"encoding/json"

	"github.com/godcong/oauth2server/model"
)

func TestParseClientUser(t *testing.T) {
	form := "client_id=cid0001&redirect_uri=https%3A%2F%2Ftest3.mana.cn%2Fauth%2Fmana_callback&response_type=code"
	v, _ := url.ParseQuery(form)
	fmt.Println(v)
}

func TestCreateUser(t *testing.T) {
	user := model.NewUser()
	user.Mobile = "123"
	user.Username = "fdsafdsa"
	user.GenerateBCryptPassword("12345")
	user.RegisterType = "cidejkfdsa"
	log.Println("regiseter save: " + user.Username + "|" + user.Mobile)
	//	model.Gorm().Create(user)
	log.Println(model.Save(user))
}

func TestParseURL(t *testing.T) {
	log.Println(os.ModeAppend)
	log.Println(os.ModePerm)

}

func TestHome(t *testing.T) {
	//out := gin.H{
	//	"sub":          "author.SubID",
	//	"nickname":     "user.Nickname",
	//	"name":         "user.Username",
	//	"phone_number": "user.Mobile",
	//	"email":        "user.Mail",
	//	"picture":      ""}
	//
	//json.NewEncoder(os.Stdout).Encode(out)

	sdata := `[{"uid":"0MLvoHrRiyfJIIGrYX65oXloEbDZ9rpr","message":"SUCCESS"},
{"uid":"hAqPzRACTpNVTgNXvAHmTTL5CWvBVeKS",  "message":"SUCCESS"},
{"uid":"gSbVT4zMo8IR1OAbNwRN8DyPTS6z1OFS",  "message":"含有敏感词"},
{"uid":"hAqPzRACTpNVTgNXvAHmTTL5CWvBVeKS",  "message":"哈哈哈"}]`

	var m []struct {
		Uid     string `json:"uid"`
		Message string `json:"message"`
	}

	//json.NewDecoder(strings.NewReader(sdata)).Decode(&m)
	json.Unmarshal([]byte(sdata), &m)
	fmt.Println(m)
}

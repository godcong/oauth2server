package oauth2server

import (
	//	"encoding/json"

	"net/http"

	"time"

	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/godcong/oauth2server/base"
	"github.com/godcong/oauth2server/model"
	"github.com/godcong/oauth2server/msg"
)

var (
	AuthorizationCodeGrant   = "authorization_code"
	PasswordCredentialsGrant = "password"
	ClientCredentialsGrant   = "client_credentials"
	RefreshToken             = "refresh_token"
)

func tokenGet(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"error": "Request must be POST",
	})
}

//func tokenPost(c *gin.Context) {
//>>>>>>> 61ccd1cfb55d69afe946a99abd5151b5ce3cad0b

func tokenPost(c *gin.Context) {
	fmt.Println("tokenPost")
	if c.Request.Form == nil {
		c.Request.ParseForm()
	}

	if e := c.Err(); e != nil {
		log.Println("err1", e.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": e.Error(),
		})
		return

	}
	if e := AccessCheck(c); e != nil {

		log.Println("err2", e.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": e.Error(),
		})
		return
	}

	form := c.Request.Form
	grant_type := form.Get("grant_type")

	switch grant_type {
	case AuthorizationCodeGrant:
		Token(c)
	case PasswordCredentialsGrant:
		PasswordToken(c)
	case RefreshToken:
		GetRefreshToken(c)
	default:
		msg.Info("err3")

		c.JSON(http.StatusBadRequest, gin.H{
			"error": ERROR_MAP[E_INVALID_GRANT],
		})
		return
	}
}

//获取Token
func Token(c *gin.Context) {

	if c.Request.Form == nil {
		c.Request.ParseForm()
	}

	form := c.Request.Form
	code := form.Get("code")
	redirect_uri := form.Get("redirect_uri")

	if code == "" || redirect_uri == "" {
		log.Println("err3")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": ERROR_MAP[E_INVALID_REQUEST],
		})
		return
	}

	j, e := GetRedis().Do("GET", code)
	if e != nil {
		log.Println("err4")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": e.Error(),
		})
		return
	}

	user := model.User{}

	base.JsonDecode(string(j.([]byte)), &user)
	if user.IsNull() == true {
		log.Println("err5")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": ERROR_MAP[E_INVALID_REQUEST],
		})
		return
	}
	client := new(model.Client)
	log.Println(c.Request.Form.Get("client_id"))
	model.Gorm().First(client, "client_user = ?", c.Request.Form.Get("client_id"))

	if client.IsNull() {
		log.Println("client is null")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": ERROR_MAP[E_INVALID_REQUEST],
		})
		return
	}

	acc, ref := GetAccessTokenAndRefreshToekn(client, &user)

	c.JSON(http.StatusOK, gin.H{
		"access_token":  acc,
		"token_type":    "bearer",
		"refresh_token": ref,
		"expires_in":    "0",
		//"expires":       time.Nanosecond.Nanoseconds(),
	})
}

//token过期，使用refresh_token重新获取
func GetRefreshToken(c *gin.Context) {

	if c.Request.Form == nil {
		c.Request.ParseForm()
	}

	form := c.Request.Form
	client_id := form.Get("client_id")
	refresh_token := form.Get("refresh_token")

	if refresh_token == "" {
		log.Println("err3")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": ERROR_MAP[E_INVALID_REQUEST],
		})
		return
	}

	authorize := new(model.Authorize)

	model.Gorm().First(authorize, "refresh_token = ?", refresh_token)

	acc, ref := AccessGenerateToken(client_id, authorize.UserID.String(), time.Nanosecond.Nanoseconds(), true)

	authorize.AccessToken = acc
	authorize.RefreshToken = ref

	GetRedis().Do("SET", acc, authorize.SubID)

	c.JSON(http.StatusOK, gin.H{
		"access_token":  acc,
		"token_type":    "bearer",
		"expires_in":    3600,
		"refresh_token": ref,
		//		"example_parameter": "example_value",
	})
}

func PasswordToken(c *gin.Context) {

	if c.Request.Form == nil {
		c.Request.ParseForm()
	}

	form := c.Request.Form
	client_id := form.Get("client_id")
	username := form.Get("username")
	passwd := form.Get("password")

	if username == "" || passwd == "" {
		msg.Info("err")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": ERROR_MAP[E_INVALID_REQUEST],
		})
		return
	}

	if scope := form.Get("scope"); scope != "" {
		GetRedis().Do("SET", client_id, scope)
	}

	user := new(model.User)
	model.Gorm().First(user, "username = ?", username)

	if user.Password != passwd {
		fmt.Println("username or passwortd is wrong!!")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": ERROR_MAP[E_INVALID_REQUEST],
		})
		return
	}

	client := new(model.Client)
	model.Gorm().First(client, "client_user = ?", client_id)
	if client.IsNull() {
		msg.Info("client is null")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": ERROR_MAP[E_INVALID_REQUEST],
		})
		return
	}
	acc, ref := GetAccessTokenAndRefreshToekn(client, user)

	c.JSON(http.StatusOK, gin.H{
		"access_token":  acc,
		"token_type":    "bearer",
		"expires_in":    3600,
		"refresh_token": ref,
	})

}

/*
validate:
client_id	Required. The client application's id.
client_secret	Required. The client application's client secret .
grant_type	Required. Must be set to authorization_code .
code	Required. The authorization code received by the authorization server.
redirect_uri	Required, if the request URI was included in the authorization request. Must be identical then.
*/
func GetAccessTokenAndRefreshToekn(client *model.Client, user *model.User) (acc, ref string) {

	author := new(model.Authorize)
	oid := model.GenerateOpenID(client.ID.String(), user.ID.String())

	author.SubID = oid
	author.ClientID = client.ID
	author.UserID = user.ID

	acc, ref = AccessGenerateToken(client.ClientUser, user.ID.String(), time.Nanosecond.Nanoseconds(), true)

	if exts, author2 := CreateIfNotExist(author); exts {
		author.RefreshToken = ref
		author.AccessToken = acc
		_ = author2
	}

	GetRedis().Do("SET", acc, author.SubID)
	return
}

func AccessCheck(c *gin.Context) error {

	if c.Request.Form == nil {
		c.Request.ParseForm()
	}

	client_id := c.Request.Form.Get("client_id")
	client_secret := c.Request.Form.Get("client_secret")
	grant_type := c.Request.Form.Get("grant_type")
	//	code := c.Request.Form.Get("code")
	//	redirect_uri := c.Request.Form.Get("redirect_uri")

	if client_id == "" || client_secret == "" || grant_type == "" {
		log.Println("missed", c.Request.Form)
		return ERROR_MAP[E_INVALID_REQUEST]
	}

	if c.Request.Method != "POST" {
		log.Println("not post")
		return ERROR_MAP[E_INVALID_REQUEST]
	}

	client := new(model.Client)
	model.Gorm().First(client, "client_user = ?", client_id)

	if client.IsNull() {
		log.Println("client is null")
		return ERROR_MAP[E_INVALID_CLIENT]
	}

	if client.GetSecret() != client_secret {
		log.Println("client secret is wrong!")
		return ERROR_MAP[E_INVALID_CLIENT]
	}

	if grant_type == "authorization_code" {
		flag := client.CheckRedirectUri(c.Request.Form.Get("redirect_uri"))
		if !flag {
			log.Println("RedirectUri is wrong！")
			return ERROR_MAP[E_INVALID_GRANT]
		}
	}

	return ERROR_MAP[E_INVALID_NONE]

}

func CreateIfNotExist(author *model.Authorize) (bool, *model.Authorize) {
	existAuthor := new(model.Authorize)
	model.Gorm().First(existAuthor, "sub_id = ?", author.SubID)
	log.Println("CreateIfNotExist ID: " + existAuthor.ID.String())
	log.Println("SubID1: " + author.SubID)
	log.Println("SubID2: " + existAuthor.SubID)
	if existAuthor.SubID == author.SubID {
		return true, existAuthor
	}
	model.Gorm().Create(author)
	return false, nil
}

func RefreshCheck() {

}

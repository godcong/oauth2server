package oauth2server

import (
	"bytes"
	"encoding/base64"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/godcong/oauth2server/model"

	"github.com/godcong/oauth2server/base"

	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
)

const TOKEN_PREFIX = "TK_"

func FinishAuthorize(c *gin.Context, u interface{}) (code string, e error) {
	if c.Request.Form == nil {
		c.Request.ParseForm()
	}
	form := c.Request.Form.Encode()

	f, e := url.ParseQuery(form)
	if e != nil {
		return "", ERROR_MAP[E_INVALID_REQUEST]
	}

	if e := CheckClientValidator(f); e != E_INVALID_NONE {
		return "", ERROR_MAP[e]
	}

	if e := c.Request.ParseForm(); e != nil {
		return "", e
	}

	rt := f.Get("response_type")
	log.Println(rt)
	if rt == "code" {
		state := f.Get("state")
		cli := f.Get("client_id")
		uri := f.Get("redirect_uri")

		if scope := f.Get("scope"); scope != "" {
			GetRedis().Do("SET", cli, scope)
		}

		code = AuthorizeGenerateToken(cli, u.(model.User).ID.String())

		res := base.StitchAddress(
			[]string{"code", code},
			[]string{"state", state},
		)

		c.Redirect(http.StatusFound, strings.Join([]string{uri, res}, "?"))
	} else if rt == "token" {
		cli := f.Get("client_id")
		uri := f.Get("redirect_uri")
		scope := f.Get("scope")
		state := f.Get("state")

		client := new(model.Client)
		model.Gorm().First(client, "client_user = ?", cli)
		if client.IsNull() {
			return
		}

		if scope != "" {
			GetRedis().Do("SET", cli, scope)
		}

		author := new(model.Authorize)
		oid := model.GenerateOpenID(cli, u.(model.User).ID.String())

		author.SubID = oid
		author.ClientID = client.ID
		author.UserID = u.(model.User).ID

		code = AuthorizeGenerateToken(cli, u.(model.User).ID.String())
		acc, _ := AccessGenerateToken(cli, u.(model.User).ID.String(), time.Nanosecond.Nanoseconds(), false)

		GetRedis().Do("SET", acc, author.SubID)
		res := base.StitchAddress(
			[]string{"access_token", acc},
			[]string{"token_type", "bearer"},
			[]string{"expires_in", "3600"},
			[]string{"scope", scope},
			[]string{"state", state},
		)
		c.Redirect(http.StatusFound, strings.Join([]string{uri, res}, "?"))
	} else {
		log.Println("Bad Request!!!")
		c.Redirect(http.StatusBadRequest, "/home")
	}
	return code, nil
}

func CheckClientValidator(form url.Values) int {
	client := new(model.Client)

	id := form.Get("client_id")
	redirectUri := form.Get("redirectUri")
	log.Println("CheckClientValidator client_id: ", id)
	log.Println("CheckClientValidator redirectUri: ", redirectUri)
	if redirectUri == "" {
		return E_UNAUTHORIZED_CLIENT
	}

	model.Gorm().First(client, "client_user = ?", id)
	if client.IsNull() {
		log.Println("authorize client is null")
		return E_UNAUTHORIZED_CLIENT
	}

	flag := client.CheckRedirectUri(redirectUri)

	//	if client.RedirectUri == redirectUri || client.RedirectUri == url.QueryEscape(redirectUri) {
	if flag || client.RedirectUri == url.QueryEscape(redirectUri) {
		log.Println("RedirectUri success")
		return E_INVALID_NONE
	}

	log.Println("first", client.RedirectUri)
	log.Println("second", redirectUri)
	log.Println("no catched")
	return E_UNAUTHORIZED_CLIENT

}

func AuthorizeGenerateToken(cid, uid string) (code string) {
	buf := bytes.NewBufferString(cid)
	buf.WriteString(uid)

	token := uuid.NewV3(uuid.NewV1(), buf.String())
	code = base64.URLEncoding.EncodeToString(token.Bytes())
	code = TOKEN_PREFIX + strings.ToUpper(strings.TrimRight(code, "="))
	return
}

func AccessGenerateToken(cid, uid string, nano int64, genRefresh bool) (access, refresh string) {
	buf := bytes.NewBufferString(cid)
	buf.WriteString(uid)
	buf.WriteString(strconv.FormatInt(nano, 10))

	access = base64.URLEncoding.EncodeToString(uuid.NewV3(uuid.NewV4(), buf.String()).Bytes())
	access = TOKEN_PREFIX + strings.ToUpper(strings.TrimRight(access, "="))
	if genRefresh {
		refresh = base64.URLEncoding.EncodeToString(uuid.NewV5(uuid.NewV4(), buf.String()).Bytes())
		refresh = TOKEN_PREFIX + strings.ToUpper(strings.TrimRight(refresh, "="))
	}
	return
}

func ResponseError(c *gin.Context, code int) {
	if err := GetError(code); err != nil {
		c.Error(err)
	}

	return
}

//验证客户端信息
//失败跳转到主页
//验证错误跳转到主页
func authorizeGet(c *gin.Context) {
	session := AutoDefaultSession(c)
	j := session.Get("LoggedUserInfo")
	user := new(model.User)
	err := base.JsonDecode(j, user)
	if j == "" && err != nil {
		log.Println("no info")
		redirectToSign(c)
		return
	}

	client := getClient(c)

	if !client.IsNull() && client.Type == 0 {
		log.Println("client authorize")
		code, e := FinishAuthorize(c, *user)
		if e == nil {
			GetRedis().Do("SET", code, j)
		} else {
			c.JSON(http.StatusOK, gin.H{
				"error": e.Error(),
			})
			return
		}
	}

	cli := getClient(c)

	if base.JsonDecode(j, user) == nil {
		c.HTML(http.StatusOK, "authorize.html", gin.H{"title": GetTitleFromClient(cli), "user": *user, "client": *cli})
		return
	}
	log.Println("default to sign")
	redirectToSign(c)
	return
}

//验证通过回跳
func authorizePost(c *gin.Context) {
	session := AutoDefaultSession(c)

	j := session.Get("LoggedUserInfo")

	user := model.User{}
	if j == "" || base.JsonDecode(j, &user) != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": "system error",
		})
		return
	}

	code, e := FinishAuthorize(c, user)
	if e == nil {
		GetRedis().Do("SET", code, j)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"error": e.Error(),
		})
	}

	return
}

func getClient(c *gin.Context) *model.Client {

	if c.Request.Form == nil {
		c.Request.ParseForm()
	}
	form := c.Request.Form
	log.Println(form)
	cid := form.Get("client_id")
	client := new(model.Client)
	model.Gorm().First(client, "client_user = ?", cid)

	return client

}

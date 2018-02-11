package oauth2server

import (
	"log"
	"net/http"
	"net/url"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/godcong/oauth2server/base"
	"github.com/godcong/oauth2server/model"
)

const DefaultClientName = "玛娜花园"

//用户登录界面
func loginGet(c *gin.Context) {

	session := AutoDefaultSession(c)
	user := model.User{}
	j := session.Get("LoggedUserInfo")
	if j != "" {
		session.Delete("LoggedUserInfo")
		e := base.JsonDecode(j, &user)
		if e == nil && user.ID.String() != "" {

		}
	}

	if c.Request.Form == nil {
		c.Request.ParseForm()
	}

	if c.Request.Form.Get("client_id") != "" {
		session.Set("Form", c.Request.Form.Encode())
	}

	client := ParseClient(c.Request.Form.Encode())

	reg := ParseURL(c, "/register")
	reg.RawQuery = c.Request.URL.RawQuery
	forg := ParseURL(c, "/forget")
	forg.RawQuery = c.Request.URL.RawQuery

	c.HTML(http.StatusOK, "login.html", gin.H{
		"title":     GetTitleFromClient(client),
		"regurl":    reg.String(),
		"forgeturl": forg.String(),
	})
	log.Println("reg:", reg.String())
	log.Println("forg:", forg.String())
	return
}

func GetTitleFromClient(client *model.Client) string {
	if client != nil && !client.IsNull() && client.ClientName != "" {
		return client.ClientName
	}
	return DefaultClientName
}

//登录提交页面
func loginPost(c *gin.Context) {
	session := AutoDefaultSession(c)
	funcWay := make(map[string]func(c *gin.Context, s *Session) bool)
	funcWay["1"] = ValidatePhone
	funcWay["0"] = ValidateUser

	way := c.DefaultPostForm("way", "0")

	if f, b := funcWay[way]; b == true {
		if !f(c, session) {
			log.Println("check failed")
			return
		}
	}

	rlt := A_LOGIN_SUCCESS
	url := ParseURL(c, "/home")
	if f := PullForm(session); f != "" {
		url.Path = "authorize"
		url.RawQuery = f
	}
	rlt.Data = map[string]string{
		"URL": url.String(),
	}
	log.Println("login success")
	c.JSON(200, rlt)

	return
}

func PullForm(s *Session) string {
	form := s.Get("Form")
	defer s.Delete("Form")
	return form

}

func ParseURL(c *gin.Context, path string) *url.URL {
	u := new(url.URL)
	u.Path = path
	return u
}

func ValidatePhone(c *gin.Context, s *Session) bool {

	mobile := c.DefaultPostForm("mobile", "")
	code := c.DefaultPostForm("code", "")
	if mobile == "" || code == "" {
		//c.Error(errors.New("username or password is wrong"))
		c.JSON(200, A_MOBILE_OR_CODE_CANNOT_NULL)
		return false
	}

	log.Println("ValidatePhone code:" + code)
	log.Println("ValidatePhone mobile:" + mobile)

	cid := ParseClientUser(c)
	log.Println("register cid: " + cid)

	user := new(model.User)

	t := model.VerifyAccountType(mobile)
	if t == model.ACCOUNT_TYPE_NONE {
		c.JSON(200, A_MOBILE_LENGTH_WRONG)
		return false
	}

	rmap := base.MessageCheck(mobile, code)
	log.Println(rmap)
	if rmap == nil {
		c.JSON(http.StatusOK, A_MESSAGE_CHECK_FAILED)
		return false
	}

	if v, b := (*rmap)["code"]; b == false || v != "0" {
		c.JSON(http.StatusOK, A_MESSAGE_CHECK_FAILED)
		return false
	}

	user.GetUser(mobile, t)
	log.Println(mobile)
	if user.IsNull() == true {
		user = model.NewUser()
		user.Mobile = mobile

		user.RegisterType = cid
		model.Save(user)
		//=======
		//		model.Gorm().Create(user)
		//		log.Println("create:", model.Gorm().Error.Error())
		//>>>>>>> 61ccd1cfb55d69afe946a99abd5151b5ce3cad0b
		user.GetUser(mobile, t)
	}
	log.Println("ValidatePhone uid: " + user.ID.String())

	user.UpdateSignInfo(base.ObtainClientIP(c.Request), time.Now().String(), "")
	j, e := base.JsonEncode(user)
	if e == nil {
		log.Println("loguser:", j)
		s.Set("LoggedUserInfo", j)
	}

	model.Save(user)
	//log.Println("save:", model.Gorm().Error.Error())
	return true
}

// return: continue flag
func ValidateUser(c *gin.Context, s *Session) bool {
	uname := c.DefaultPostForm("username", "")
	password := c.DefaultPostForm("password", "")

	if uname == "" || password == "" {
		//c.Error(errors.New("username or password is wrong"))
		c.JSON(200, A_NAME_OR_PASSWORD_CANNOT_NULL)
		return false
	}

	user := model.User{}
	t := model.VerifyAccountType(uname)
	if t == model.ACCOUNT_TYPE_NONE {
		c.JSON(200, A_NAME_OR_PASSWORD_WRONG)
		return false
	}

	user.GetUser(uname, t)

	if t == model.ACCOUNT_TYPE_MOBILE && user.Password == "" {
		//跳转到完善用户信息页面
		user.UpdateSignInfo(base.ObtainClientIP(c.Request), time.Now().String(), "")
		j, e := base.JsonEncode(user)
		if e == nil {
			s.Set("CompleteUserInfo", j)
		}

		rlt := A_LOGIN_SUCCESS
		url := ParseURL(c, "/completeuserinfo")
		rlt.Data = map[string]string{
			"URL": url.String(),
		}
		c.JSON(200, rlt)
		return false
	}

	if !user.VerifyPassword(password) {
		c.JSON(200, A_NAME_OR_PASSWORD_WRONG)
		return false
	}

	user.UpdateSignInfo(base.ObtainClientIP(c.Request), time.Now().String(), "")
	j, e := base.JsonEncode(user)
	if e == nil {
		s.Set("LoggedUserInfo", j)
	}
	model.Save(user)
	return true
}

func registerGet(c *gin.Context) {
	session := AutoDefaultSession(c)
	if c.Request.Form == nil {
		c.Request.ParseForm()
	}
	log.Println(c.Request.Form)
	if c.Request.Form.Get("client_id") != "" {
		session.Set("Form", c.Request.Form.Encode())
	}
	cli := getClient(c)

	c.HTML(http.StatusOK, "register.html", gin.H{"title": GetTitleFromClient(cli)})
}

func registerPost(c *gin.Context) {
	session := AutoDefaultSession(c)
	if c.Request.Form == nil {
		c.Request.ParseForm()
	}
	log.Println(c.Request.Form)

	uname := c.DefaultPostForm("username", "")
	mobile := c.DefaultPostForm("mobile", "")
	password := c.DefaultPostForm("password", "")
	vpassword := c.DefaultPostForm("vpassword", "")
	code := c.DefaultPostForm("code", "")

	cid := ParseClientUser(c)
	log.Println("register cid: " + cid)

	log.Println("username:", uname, " mobile:", mobile, " password:", password, " code:", code) ///
	user := model.NewUser()
	if uname == "" || password == "" || vpassword == "" {
		c.JSON(http.StatusOK, A_MUST_NOT_BE_NULL)
		return
	}

	if len(uname) < 6 {
		c.JSON(http.StatusOK, A_NAME_LENGTH_WRONG)
		return
	}

	if b, e := model.VerifyUsername(uname); b == false || e != nil {
		c.JSON(http.StatusOK, A_NAME_FORMAT_WRONG)
		return
	}

	user.GetUserByUsername(uname)
	if !user.IsNull() {
		c.JSON(http.StatusOK, A_NAME_IS_EXISTS)
		return
	}

	//if b, e := model.VerifyMobile(mobile); b == false || e != nil {
	//	c.JSON(http.StatusOK, A_MOBILE_LENGTH_WRONG)
	//	return
	//}

	//user.GetUserByMobile(mobile)
	//if !user.IsNull() {
	//	c.JSON(http.StatusOK, A_MOBILE_IS_EXISTS)
	//	return
	//}

	//rmap := base.MessageCheck(mobile, code)
	//
	//if rmap == nil {
	//	c.JSON(http.StatusOK, A_MESSAGE_CHECK_FAILED)
	//	return
	//}

	//if v, b := (*rmap)["code"]; b == false || v != "0" {
	//	c.JSON(http.StatusOK, A_MESSAGE_CHECK_FAILED)
	//	return
	//}
	log.Println("register uid: " + user.ID.String())
	if user.IsNull() {
		user.Mobile = mobile
		user.Username = uname
		user.GenerateBCryptPassword(password)
		user.RegisterType = cid
		log.Println("regiseter save: " + user.Username + "|" + user.Mobile)
		model.Gorm().Create(user)
	}
	j, e := base.JsonEncode(user)
	url := ParseURL(c, "/login")
	if e == nil {
		if f := PullForm(session); f != "" {
			url.Path = "authorize"
			url.RawQuery = f
		}
		session.Set("LoggedUserInfo", j)
	}

	rlt := A_REGISTER_SUCCESS

	rlt.Data = map[string]string{
		"URL": url.String(),
	}

	c.JSON(http.StatusOK, rlt)
	return

}

func ParseClient(form string) *model.Client {
	if form == "" {
		return nil
	}

	client := model.NewClient()
	if fmap, e := url.ParseQuery(form); e == nil {

		if cid, b := fmap["client_id"]; b == true {
			model.Gorm().First(client, "client_user = ?", cid)
			if !client.IsNull() {
				return client
			}
		}
	}
	return nil
}

func ParseClientUser(c *gin.Context) string {
	session := AutoDefaultSession(c)
	form := session.Get("Form")

	log.Println("ParseClientUser form: " + form)
	if form == "" {
		return ""
	}

	client := model.NewClient()
	if fmap, e := url.ParseQuery(form); e == nil {

		if cid, b := fmap["client_id"]; b == true {
			model.Gorm().First(client, "client_user = ?", cid)
			if !client.IsNull() {
				return client.ClientUser
			}
		}
	}
	return ""
}

func registerPhoneSend(c *gin.Context) {

	if c.Request.Form == nil {
		c.Request.ParseForm()
	}
	log.Println(c.Request.Form)
	mobile := c.DefaultPostForm("mobile", "")
	log.Println("mobile", mobile, reflect.TypeOf(mobile))
	if b, e := model.VerifyMobile(mobile); b == false || e != nil {
		log.Println(b, e, mobile)
		c.JSON(http.StatusOK, A_MOBILE_LENGTH_WRONG)
		return
	}

	rmap := base.RegisterSend(mobile)
	if rmap == nil {

		c.JSON(http.StatusOK, A_MESSAGE_SEND_FAILED)
		return
	}

	if v, b := (*rmap)["code"]; !b || v == "1" {
		c.JSON(http.StatusOK, A_MESSAGE_SEND_FAILED)
		return
	}

	c.JSON(http.StatusOK, A_MESSAGE_SEND_SUCCESS)
	return
}

func registerPhoneCheck(c *gin.Context) {

	if c.Request.Form == nil {
		c.Request.ParseForm()
	}
	mobile := c.DefaultPostForm("mobile", "")
	code := c.DefaultPostForm("code", "")

	if b, e := model.VerifyMobile(mobile); b == false || e != nil {
		c.JSON(http.StatusOK, A_MOBILE_LENGTH_WRONG)
	}

	if code == "" {
		c.JSON(http.StatusOK, A_MESSAGE_WRONG)
		return
	}

	rmap := base.MessageCheck(mobile, code)

	if rmap == nil {
		c.JSON(http.StatusOK, A_MESSAGE_CHECK_FAILED)
		return
	}
	c.JSON(http.StatusOK, A_MESSAGE_CHECK_SUCCESS)
	return
}

func forgetGet(c *gin.Context) {
	cli := getClient(c)
	c.HTML(http.StatusOK, "forget.html", gin.H{"title": GetTitleFromClient(cli)})
}

func resetGet(c *gin.Context) {
	cli := getClient(c)
	c.HTML(http.StatusOK, "reset.html", gin.H{"title": GetTitleFromClient(cli)})
}

func forgetPost(c *gin.Context) {

	session := AutoDefaultSession(c)
	if c.Request.Form == nil {
		c.Request.ParseForm()
	}

	if uid := session.Get("ForgetUser"); uid == "" {
		log.Println(uid)

		mobile := c.DefaultPostForm("mobile", "")
		uname := c.DefaultPostForm("username", "")
		mail := c.DefaultPostForm("mail", "")
		code := c.DefaultPostForm("code", "")

		log.Println("mobile:", mobile, " uname:", uname, " mail:", mail, " code:", code)

		if mobile == "" && uname == "" && mail == "" && code == "" {
			c.JSON(http.StatusOK, A_FORGET_NEED_IS_NULL)
			return
		}

		user := model.User{}

		if mobile != "" {
			if b, e := model.VerifyMobile(mobile); b == true && e == nil {
				user.GetUser(mobile, model.ACCOUNT_TYPE_MOBILE)
			}
		} else if uname != "" {
			if b, e := model.VerifyUsername(mobile); b == true && e == nil {
				user.GetUser(mobile, model.ACCOUNT_TYPE_UNAME)
			}
		} else {
		}

		if user.Mobile == "" {
			c.JSON(http.StatusOK, A_FORGET_ACCOUNT_WRONG)
			return

		}

		rmap := base.MessageCheck(mobile, code)

		if rmap == nil {
			c.JSON(http.StatusOK, A_MESSAGE_CHECK_FAILED)
			return
		}

		if v, b := (*rmap)["code"]; b == false || v != "0" {
			c.JSON(http.StatusOK, A_MESSAGE_CHECK_FAILED)
			return
		}

		session.Set("ForgetUser", user.ID.String())
		c.JSON(http.StatusOK, A_FORGET_CHECK_SUCCESS)
	} else {

		log.Println(uid)

		rf := c.DefaultPostForm("type", "")
		pass := c.DefaultPostForm("password", "")
		vpass := c.DefaultPostForm("vpassword", "")

		log.Println("type:", rf, " pass:", pass, " vpass:", vpass)
		if rf != "reset" {
			session.Delete("ForgetUser")
			c.JSON(http.StatusOK, A_FORGET_SYSTEM_ERROR)
			return
		}

		if pass == "" || vpass == "" {
			c.JSON(http.StatusOK, A_FORGET_PASSWORD_CANNOT_NULL)
			return
		}

		if pass != vpass {
			c.JSON(http.StatusOK, A_FORGET_PASSWORD_FILL_OUT)
			return

		}

		user := model.User{}

		model.FirstById(&user, uid)

		if user.ID.String() == "" {
			c.JSON(http.StatusOK, A_FORGET_ACCOUNT_WRONG)
			return
		}
		//		user.UpdateSignInfo(base.ObtainClientIP(c.Request), time.Now().String(), "")
		u := model.NewUser()
		//		user.GenerateBCryptPassword(pass)
		u.UpdateSignInfo(base.ObtainClientIP(c.Request), time.Now().String(), "")
		u.GenerateBCryptPassword(pass)

		log.Println("更改后的密码：", u.Password)
		//	model.Save(&user)
		//		model.Gorm().Model(user).Update(u)
		model.Update(&user, u)
		session.Delete("ForgetUser")
		c.JSON(http.StatusOK, A_FORGET_CHANGE_SUCCESS)
		return
	}

	return

}

func changeGet(c *gin.Context) {
	cli := getClient(c)
	c.HTML(http.StatusOK, "change.html", gin.H{"title": GetTitleFromClient(cli)})
}

func changePost(c *gin.Context) {

	session := AutoDefaultSession(c)
	if c.Request.Form == nil {
		c.Request.ParseForm()
	}
	form := c.Request.Form
	//	opasswd(原密码),npasswd(新密码),vpasswd(确认新密码)
	opasswd := form.Get("opasswd")
	npasswd := form.Get("npasswd")
	vpasswd := form.Get("vpasswd")

	log.Println("opasswd:", opasswd, "  npasswd:", npasswd, "  vpasswd:", vpasswd)
	if opasswd == "" || npasswd == "" || vpasswd == "" {
		c.JSON(http.StatusOK, A_CHANGE_NEED_IS_NULL)
		return
	}

	j := session.Get("LoggedUserInfo")
	user := new(model.User)
	err := base.JsonDecode(j, user)

	if err != nil {
		log.Println("NO info")
		return
	}

	if !user.VerifyPassword(opasswd) {
		c.JSON(http.StatusOK, A_CHANGE_OPASSWORD_WRONG)
		return
	}

	if npasswd != vpasswd {
		c.JSON(http.StatusOK, A_CHANGE_NPASSWORD_FILL_OUT)
		return
	}

	u := model.NewUser()
	u.GenerateBCryptPassword(npasswd)

	//	u.Password = npasswd
	//	model.Gorm().Model(&user).Update(u)
	model.Update(user, u)
	log.Println("密码修改成功！！！")
	c.JSON(http.StatusOK, A_CHANGE_SUCCESS)
	return
}

func logout(c *gin.Context) {

	if c.Request.Form == nil {
		c.Request.ParseForm()
	}

	log.Println("Request:", c.Request.Form)
	session := AutoDefaultSession(c)
	session.Delete("LoggedUserInfo")
	redirectToSign(c)
}

//跳转到登陆页面，并保留uri参数
func redirectToSign(c *gin.Context) {

	log.Println("request:", c.Request.Form)
	u := RawUrl("/login", c.Request.Form)

	log.Println("sign: " + u.String())
	c.Writer.Header().Set("Location", u.String())
	c.Writer.WriteHeader(http.StatusFound)
	return
}

////跳转到完善用户信息页面
//func redirectToCompleteUserInfo(c *gin.Context) {
//	u := new(url.URL)
//	u.Path = "/completeuserinfo"

//	c.Writer.Header().Set("Location", u.String())
//	c.Writer.WriteHeader(http.StatusFound)
//	return
//}

func completeUserInfoGet(c *gin.Context) {
	c.HTML(http.StatusOK, "completeuserinfo.html", gin.H{})
}

//完善用户信息
func completeUserInfoPost(c *gin.Context) {

	session := AutoDefaultSession(c)
	if c.Request.Form == nil {
		c.Request.ParseForm()
	}

	uname := c.DefaultPostForm("username", "")
	password := c.DefaultPostForm("password", "")
	vpassword := c.DefaultPostForm("vpassword", "")

	code := c.DefaultPostForm("code", "")

	log.Println("username:", uname, " password:", password, " vpassword:", vpassword, " code:", code)

	cid := ParseClientUser(c)
	log.Println("register cid: " + cid)

	if uname == "" || password == "" || vpassword == "" {
		c.JSON(http.StatusOK, A_MUST_NOT_BE_NULL)
		return
	}

	if len(uname) < 6 {
		c.JSON(http.StatusOK, A_NAME_LENGTH_WRONG)
		return
	}

	if b, e := model.VerifyUsername(uname); b == false || e != nil {
		c.JSON(http.StatusOK, A_NAME_FORMAT_WRONG)
		return
	}

	if password != vpassword {
		c.JSON(http.StatusOK, A_PASSWORD_COMPARE_WRONG)
		return
	}

	user := model.NewUser()
	user.GetUserByUsername(uname)
	if !user.IsNull() {
		c.JSON(http.StatusOK, A_NAME_IS_EXISTS)
		return
	}

	u := model.NewUser()
	j := session.Get("CompleteUserInfo")
	var mobile string
	if j != "" {
		e := base.JsonDecode(j, &u)
		if e == nil {
			mobile = u.Mobile
			//通过手机号找到用户
			u.GetUserByMobile(mobile)
		}
	}

	//rmap := base.MessageCheck(mobile, code)
	//
	//if rmap == nil {
	//	c.JSON(http.StatusOK, A_MESSAGE_CHECK_FAILED)
	//	return
	//}
	//
	//if v, b := (*rmap)["code"]; b == false || v != "0" {
	//	c.JSON(http.StatusOK, A_MESSAGE_CHECK_FAILED)
	//	return
	//}

	if user.IsNull() {
		//更新信息
		user1 := model.NewUser()
		user1.Username = uname
		user1.GenerateBCryptPassword(password)
		user1.RegisterType = cid
		log.Println("regiseter save: " + user1.Username + "|" + user1.Mobile)
		//		model.Gorm().Model(&u).Update(user1)
		model.Update(u, user1)
	}
	c.JSON(http.StatusOK, A_COMPLETE_USERINFO_SUCCESS)
	return

}
func RawUrl(s string, form url.Values) *url.URL {
	u := new(url.URL)
	u.Path = s
	log.Println("form.Encode:", form.Encode())
	u.RawQuery = form.Encode()
	return u

}

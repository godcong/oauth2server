package oauth2server

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/godcong/oauth2server/base"
)

func SignCheck(c *gin.Context) {
	log.Println("SignCheck")
	session := AutoDefaultSession(c)
	user := session.Get("LoggedUserInfo")
	if user == "" {
		if c.Request.Form == nil {
			c.Request.ParseForm()
		}
		if c.Request.Form.Get("client_id") != "" {
			session.Set("Form", c.Request.Form.Encode())
		}

		redirectToSign(c)
		c.Abort()
		return
	}
	c.Next()
}

func VisitLog(c *gin.Context) {
	base.Println("VisitLog URL: ", c.Request.URL.RawQuery)
	c.Next()
}

package oauth2server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (gs *Oauth2Server) Router() {
	router(gs)
}

func router(gs *Oauth2Server) {
	homeMain := Home(gs.router)

	gs.router.LoadHTMLGlob("templates/*")

	homeMain.Any("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "login")
	})
	homeMain.GET("/login", loginGet)
	homeMain.POST("/login", loginPost)
	homeMain.GET("/register", registerGet)
	homeMain.POST("/register", registerPost)
	homeMain.GET("/forget", forgetGet)
	homeMain.POST("/forget", forgetPost)
	homeMain.GET("/reset", resetGet)

	homeMain.GET("/completeuserinfo", completeUserInfoGet)
	homeMain.POST("/completeuserinfo", completeUserInfoPost)

	homeMain.GET("/protocol", func(c *gin.Context) {
		cli := getClient(c)
		c.HTML(http.StatusOK, "protocol.html", gin.H{"title": GetTitleFromClient(cli)})
	})
	homeMain.GET("/token", tokenGet)
	homeMain.POST("/token", tokenPost)

	homeMain.GET("/userinfo", userinfoGet)
	homeMain.POST("/msg/send", registerPhoneSend)
	homeMain.POST("/msg/check", registerPhoneCheck)

	oauth := homeMain.Group("/")
	oauth.Use(SignCheck)
	oauth.GET("/logout", logout)
	oauth.GET("/authorize", authorizeGet)
	oauth.POST("/authorize", authorizePost)
	oauth.GET("/change", changeGet)
	oauth.POST("/change", changePost)
	oauth.GET("/home", home)

}

func Home(e *gin.Engine) *gin.RouterGroup {
	homeMain := e.Group("/")
	homeMain.Use(VisitLog)
	return homeMain
}

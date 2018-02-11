package oauth2server

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type Session struct {
	session sessions.Session
	auto    bool
}

var option *sessions.Options

func init() {
	option = new(sessions.Options)
	//设置session的过期时间
	option.MaxAge = 0
}

func AutoDefaultSession(c *gin.Context) *Session {
	session := new(Session)
	session.session = sessions.Default(c)
	return session
}

func (s *Session) Delete(key string) error {
	s.session.Delete(key)
	return s.session.Save()
}

func (s *Session) Set(key, val string) error {
	s.session.Set(key, val)
	s.session.Options(*option)
	return s.session.Save()
}

func (s *Session) Get(key string) string {
	switch r := s.session.Get(key); r.(type) {
	case string:
		return r.(string)
	default:
		return ""
	}
}

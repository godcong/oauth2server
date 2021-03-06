package oauth2server

import (
	"log"
	"strings"

	"github.com/garyburd/redigo/redis"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/godcong/oauth2server/config"
	"github.com/godcong/oauth2server/model"
	"github.com/jinzhu/gorm"

	"github.com/garyburd/redigo/redisx"
)

type Oauth2Server struct {
	store  *sessions.RedisStore
	router *gin.Engine
	db     *gorm.DB
	rds    *redisx.ConnMux
}

var (
	//oaserver *OAuthServer
	gs *Oauth2Server
)

func init() {
	gs = new(Oauth2Server)

	gs.store = NewRedisStore()

	gs.router = gin.Default()

	gs.router.Use(sessions.Sessions("gautu", *gs.store))
	gs.router.Static("/static", "static")
	gs.router.StaticFile("/favicon.ico", "favicon.ico")
	gs.db = model.Gorm()
	gs.rds = NewRedis()
}

type defaultRedis struct {
	Address  string
	Port     string
	Password string
	DB       string
	User     string
}

func NewRedis() *redisx.ConnMux {

	dr := defaultRedis{
		Address:  "localhost",
		Port:     "6379",
		Password: "",
		DB:       "1",
		User:     "x",
	}
	rds := config.GetSub("redis")
	dr.Address = rds.GetStringWithDefault("addr", "localhost")
	dr.Port = rds.GetStringWithDefault("port", "6379")
	dr.Password = rds.GetStringWithDefault("password", "")
	dr.DB = rds.GetStringWithDefault("db", "1")
	dr.User = rds.GetStringWithDefault("user", "")

	//addr := strings.Join([]string{dr.Address, dr.Port}, ":")

	//addr := fmt.Sprintf("redis://%s:%s@%s:%s/%s", dr.User, dr.Password, dr.Address, dr.Port, dr.DB)
	//op := redis.DialPassword(dr.Password)
	//redis.DialNetDial()
	//c, err := redis.Dial("tcp", ":6379", op)
	//log.Panicln(addr)
	addr := strings.Join([]string{dr.Address, dr.Port}, ":")
	//c, err := redis.DialURL(addr)
	//if err != nil {
	//	panic(err)
	//}
	c, err := redis.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}
	if dr.Password != "" {
		if _, err := c.Do("AUTH", dr.Password); err != nil {
			c.Close()
			panic(err)
		}
	}

	_, err = c.Do("SELECT", dr.DB)
	if err != nil {
		c.Close()
		panic(err)
	}
	cmux := redisx.NewConnMux(c)

	return cmux
}

func GetRedis() redis.Conn {
	if gs.rds == nil || ping() != nil {
		log.Println("Could not connect to Redis!")
		gs.rds = NewRedis()
	}
	return gs.rds.Get()
}

func ping() error {
	_, err := gs.rds.Get().Do("PING")
	return err
}

func DefaultDB() *gorm.DB {
	return gs.db
}

func (gs *Oauth2Server) GetDB() *gorm.DB {
	return gs.db
}

func DefaultRedisStrore() *sessions.RedisStore {
	return gs.store
}
func (gs *Oauth2Server) GetRedisStrore() *sessions.RedisStore {
	return gs.store
}

func DefaultRouter() *gin.Engine {
	return gs.router
}
func (gs *Oauth2Server) GetRouter() *gin.Engine {
	return gs.router
}

func DefaultOauth2Server() *Oauth2Server {
	return gs
}

func NewRedisStore() *sessions.RedisStore {
	dr := defaultRedis{
		Address:  "localhost",
		Port:     "6379",
		Password: "",
		DB:       "1",
		User:     "x",
	}

	rds := config.GetSub("redis")
	dr.Address = rds.GetStringWithDefault("addr", "localhost")
	dr.Port = rds.GetStringWithDefault("port", "6379")
	dr.Password = rds.GetStringWithDefault("password", "")
	dr.DB = rds.GetStringWithDefault("db", "1")
	dr.User = rds.GetStringWithDefault("user", "")

	addr := strings.Join([]string{dr.Address, dr.Port}, ":")
	store, _ := sessions.NewRedisStoreWithDB(10, "tcp", addr, dr.Password, dr.DB, []byte("secret"))
	return &store
}

func (gs *Oauth2Server) Run(addr string) {
	gs.router.Run(addr)
}

func Start() {
	gs.Router()

	gs.Run(serverAddr())
}

func serverAddr() (r string) {
	addr := config.GetSub("system").GetStringWithDefault("addr", "")
	port := config.GetSub("system").GetStringWithDefault("port", "8080")

	r = strings.Join([]string{
		addr,
		port,
	},
		":")
	return
}

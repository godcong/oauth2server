package model

import (
	"fmt"
	"log"
	"net/url"

	"github.com/godcong/oauth2server/config"

	"github.com/jinzhu/gorm"
)

type Database struct {
	Type    string
	Address string
	Port    string
	Name    string
}

type DatabaseLoginInfo struct {
	User     string
	Password string
}

type DatabaseConfiguration struct {
	Database
	DatabaseLoginInfo
	param string
	loc   string
}

var (
	defaultConfig *DatabaseConfiguration
	db            *gorm.DB
)

var defaultConfigure = &DatabaseConfiguration{
	Database: Database{
		Type:    "mysql",
		Address: "localhost",
		Port:    "3306",
		Name:    "gautu",
	},
	DatabaseLoginInfo: DatabaseLoginInfo{
		User:     "root",
		Password: "123456",
	},

	param: "?",
	loc:   url.QueryEscape("Asia/Shanghai"),
}

func init() {
	//var e error
	//defaultConfig = defaultConfigure
	//if m, e := configo.Get(defaultConfigure.Type); e == nil {
	//	defaultConfig = &DatabaseConfiguration{
	//		Database: Database{
	//			Type:    m.MustGet("db", defaultConfigure.Type),
	//			Address: m.MustGet("addr", defaultConfigure.Address),
	//			Port:    m.MustGet("port", defaultConfigure.Port),
	//			Name:    m.MustGet("dbname", defaultConfigure.Name),
	//		},
	//		DatabaseLoginInfo: DatabaseLoginInfo{
	//			User:     m.MustGet("user", defaultConfigure.User),
	//			Password: m.MustGet("password", defaultConfigure.Password),
	//		},
	//		param: m.MustGet("param", defaultConfigure.param),
	//		loc:   m.MustGet("loc", defaultConfigure.loc),
	//	}
	//
	//}
	//c := config.DefaultConfig()

	//conn := connectMysql(c)
	//base.Println(conn)
	//db = NewEngine()
}

func DatabaseType() string {
	return config.GetSub("database").GetStringWithDefault("name", "mysql")
}

func OnExit() {
	db.Close()
}

func NewEngine() *gorm.DB {
	var e error
	conn := ConnectMysql()

	db, e = gorm.Open("mysql", conn)
	if e != nil {
		log.Panic(conn)
	}
	return db

}

func ConnectMysql() string {
	db := config.GetSub("database")

	log.Println(db.Get("username"))
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s%sloc=%s&charset=utf8&parseTime=true",
		db.GetStringWithDefault("username", "root"),
		db.GetStringWithDefault("password", "123456"),
		db.GetStringWithDefault("addr", "localhost"),
		db.GetStringWithDefault("port", "3306"),
		db.GetStringWithDefault("schema", "oauth2"),
		db.GetStringWithDefault("param", "?"),
		url.QueryEscape(db.GetStringWithDefault("local", "Asia/Shanghai")))
}

func FirstById(v interface{}, id string) {
	db.Where("id = ?", id).First(v)
}

func FirstWhere(v interface{}, query string, args ...interface{}) {
	db.Where(query, args).First(v)
}

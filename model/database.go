package model

import (
	"database/sql"
	"fmt"
	"net/url"

	"gopkg.in/configo.v2"

	"log"

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
	config *DatabaseConfiguration
	db     *gorm.DB
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
	var e error
	config = defaultConfigure
	if m, e := configo.Get(defaultConfigure.Type); e == nil {
		config = &DatabaseConfiguration{
			Database: Database{
				Type:    m.MustGet("db", defaultConfigure.Type),
				Address: m.MustGet("addr", defaultConfigure.Address),
				Port:    m.MustGet("port", defaultConfigure.Port),
				Name:    m.MustGet("dbname", defaultConfigure.Name),
			},
			DatabaseLoginInfo: DatabaseLoginInfo{
				User:     m.MustGet("user", defaultConfigure.User),
				Password: m.MustGet("password", defaultConfigure.Password),
			},
			param: m.MustGet("param", defaultConfigure.param),
			loc:   m.MustGet("loc", defaultConfigure.loc),
		}

	}

	conn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s%sloc=%s&charset=utf8&parseTime=true",
		config.User, config.Password, config.Address, config.Port, config.Name, config.param, config.loc)
	log.Println(conn)
	//db.engine, e = xorm.NewEngine(db.dbtype, conn)
	db, e = gorm.Open(config.Type, conn)
	if e != nil {
		panic(e)
	}

}

func CreateDatabase() {

	db, err := sql.Open("mysql", "admin:admin@tcp(127.0.0.1:3306)/")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE DATABASE " + name)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("USE " + name)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("CREATE TABLE example ( id integer, data varchar(32) )")
	if err != nil {
		panic(err)
	}
}

func OnExit() {
	db.Close()
}

func NewEngine() *gorm.DB {

	var e error
	conn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s%sloc=%s&charset=utf8&parseTime=true",
		config.User, config.Password, config.Address, config.Port, config.Name, config.param, config.loc)

	//db.engine, e = xorm.NewEngine(db.dbtype, conn)
	db, e = gorm.Open(config.Type, conn)
	if e != nil {
		panic(e)
	}
	return db

}

func connectSql() {

}

func FirstById(v interface{}, id string) {
	db.Where("id = ?", id).First(v)
}

func FirstWhere(v interface{}, query string, args ...interface{}) {
	db.Where(query, args).First(v)
}

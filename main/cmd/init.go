package main

import (
	"database/sql"
	"fmt"
	"net/url"

	"github.com/godcong/oauth2server/config"
	"github.com/godcong/oauth2server/model"
)

func Init() {
	c := config.DefaultConfig()
	CreateDatabase(c)
	Migrate()

}

func ConnectInit(config *config.Config) string {
	db := config.GetSub("database")
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s%sloc=%s&charset=utf8&parseTime=true",
		db.GetStringWithDefault("username", "root"),
		db.GetStringWithDefault("password", "123456"),
		db.GetStringWithDefault("addr", "localhost"),
		db.GetStringWithDefault("port", "3306"),
		db.GetStringWithDefault("", ""),
		db.GetStringWithDefault("param", "?"),
		url.QueryEscape(db.GetStringWithDefault("local", "Asia/Shanghai")))
}

func CreateDatabase(config *config.Config) {

	db, err := sql.Open(model.DatabaseType(), ConnectInit(config))
	if err != nil {
		panic(err)
	}
	defer db.Close()
	//CREATE {DATABASE | SCHEMA} [IF NOT EXISTS] db_name
	//base.Print(config.GetSub("database").GetStringWithDefault("schema", "oauth2"))
	//_, err = db.Exec("CREATE SCHEMA " + "testName" + " DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci")
	_, err = db.Exec("CREATE SCHEMA " + config.GetSub("database").GetStringWithDefault("schema", "oauth2") + " DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci")
	if err != nil {
		//do nothing
	}
}

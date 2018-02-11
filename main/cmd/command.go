package main

import (
	"fmt"

	"os"
	"strconv"

	"time"

	"log"

	"net/url"

	"encoding/json"

	"github.com/godcong/oauth2server/base"
	"github.com/godcong/oauth2server/model"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"gopkg.in/configo.v2"
)

const MAX = 1000

func main() {
	args := os.Args
	fmt.Println(os.Args[0])
	if len(args) > 1 {
		if args[1] == "help" {
			Help()
		}
		if args[1] == "migrate" {
			Migrate()
		}
		if args[1] == "random" {
			num := "10"
			if len(args) > 2 {
				num = args[2]
			}
			RandomUser(num)
		}
		if args[1] == "client" {
			url := ""
			if len(args) > 2 {
				url = args[2]
			}
			ClientMake(url)
		}
		if args[1] == "user" {
			uname := string(base.GenerateRandomString(6, base.T_RAND_LOWER))
			pass := "123456"
			if len(args) > 3 {
				uname = args[2]
				pass = args[3]
			} else if len(args) > 2 {
				uname = args[2]
			}

			UserMake(uname, pass)
		}
		if args[1] == "trans" {
			//1:from 2:to
			from := "mysql2"
			to := "mysql"
			if len(args) > 1 {
				from = args[2]
			}
			if len(args) == 3 {
				to = args[3]
			}
			conF, err1 := configo.Get(from)

			if err1 != nil {
				panic(err1)
			}
			dbf := connectDb(conF)
			conT, err2 := configo.Get(to)
			if err2 != nil {
				panic(err2)
			}
			dbt := connectDb(conT)

			Transfer(dbf, dbt)
		}
	}

}
func Help() {
	fmt.Println(
		"migrate					sync the model to database\r\n" +
		"random [num]				create user of numbers(num=10)\r\n" +
		"client [url]				create an client with redirecturi(url='')\r\n" +
		"user [name] [pass]			create an user with username and password(pass=123456)\r\n")
}

func ClientMake(url string) {

	client := model.NewClient()
	id := string(base.GenerateRandomString(12, base.T_RAND_LOWERNUM))
	client.ClientUser = model.C_PREFIX_WEB + id
	client.Secret = string(base.GenerateRandomString(128))

	uri, err := json.Marshal([]string{url})
	if err == nil {
		client.RedirectUri = string(uri)
	}

	fmt.Println(client)
	db := model.Gorm().Save(client)
	fmt.Println(db.GetErrors())

}

func UserMake(name, pass string) {
	user := model.NewUser()
	user.GenerateBCryptPassword(pass)
	user.Username = name
	user2 := new(model.User)
	model.Gorm().First(user2, "username = ?", user.Username)
	if user2.IsNull() == true {
		fmt.Println(user.Username, pass)
		db := model.Save(user)
		fmt.Println(db.GetErrors())
	}

}

func Migrate() bool {
	m := model.GetModels()
	for k, v := range *m {
		fmt.Println("Migrate:", k)
		model.Gorm().AutoMigrate(v)
	}
	return true
}

func Random(args ...interface{}) bool {
	for _, v := range args {
		model.Gorm().Create(v)
	}
	return true
}

func RandomUser(num string) {

	n, e := strconv.Atoi(num)
	if e == nil {
		for ; n > 0 && n < MAX; n-- {
			user := model.NewUser()
			user.GenerateBCryptPassword("123456")
			user.Username = string(base.GenerateRandomString(6, base.T_RAND_LOWER))
			time.Sleep(1)
			user.Nickname = string(base.GenerateRandomString(12, base.T_RAND_UPPER))
			user.Mail = user.Nickname + "@mana.com"
			user.Mobile = "130" + string(base.GenerateRandomString(8, base.T_RAND_NUM))
			Random(&user)
		}
	}

}

type ManaUser struct {
	ID          uuid.UUID `gorm:"primary_key;type:varchar(36)"`
	Username    string    `gorm:"column:username"`
	Nickname    string    `gorm:"-"`
	Mobile      string    `gorm:"column:mobile"`
	Mail        string    `gorm:"column:email"`
	Password    string    `gorm:"column:password"`
	Salt        string    `gorm:"column:salt"`
	Status      string    `gorm:"column:status"`
	AccountType string    `gorm:"column:account_type"`
	ErrorTimes  int       `gorm:"column:error_times"`
	CreatedAt   int
	UpdatedAt   int
	DeletedAt   int `sql:"index"`
}

func Transfer(dbf, dbt *gorm.DB) {
	model.CFlag = false
	counts := tableCounts(dbf, "mana_users")

	lens := 10
	errs := 0
	for i := 0; i <= counts; i += lens {
		manausers := getManaUsers(dbf, i, lens)
		for _, v := range *manausers {
			if setGautuUser(dbt, v) != nil {
				errs += 1
			}
		}
	}
	counts2 := tableCounts(dbt, "users")

	log.Println("transfer", counts, "to", counts2, " errors: ", errs)
}

func setGautuUser(dto *gorm.DB, user ManaUser) error {
	u := new(model.User)
	log.Println("id:", user.ID.String())
	u.ID = user.ID
	u.Password = user.Password
	u.Salt = user.Salt
	u.ErrorTimes = 0
	u.Mail = user.Mail
	u.Mobile = user.Mobile
	u.Nickname = ""
	u.Username = ""
	u.AccountType = 10
	u.CreatedAt = parseTimestamp(user.CreatedAt)
	u.UpdatedAt = parseTimestamp(user.UpdatedAt)
	e := dto.Create(u).Error
	return e

}

func parseTimestamp(t int) time.Time {
	i, err := strconv.ParseInt(strconv.Itoa(t), 10, 64)
	if err != nil {
		return time.Now()
	}
	tm := time.Unix(i, 0)
	return tm
}

func getManaUsers(db *gorm.DB, off, lim int) *[]ManaUser {
	mu := new([]ManaUser)
	db.Order("created_at desc").Table("mana_users").Offset(off).Limit(lim).Find(mu)

	return mu
}

func tableCounts(db *gorm.DB, s string) (count int) {
	db.Order("created_at desc").Table(s).Count(&count)
	return
}

func connectDb(config *configo.Property) *gorm.DB {
	var e error
	dtyp := config.MustGet("db", "mysql")
	addr := config.MustGet("addr", "localhost")
	port := config.MustGet("port", "3306")
	name := config.MustGet("dbname", "gautu")
	user := config.MustGet("user", "root")
	pass := config.MustGet("password", "123456")
	param := config.MustGet("param", "?")
	loc := config.MustGet("loc", url.QueryEscape("Asia/Shanghai"))

	conn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s%sloc=%s&charset=utf8&parseTime=true",
		user, pass, addr, port, name, param, loc)
	log.Println(conn)
	//db.engine, e = xorm.NewEngine(db.dbtype, conn)
	db, e := gorm.Open(dtyp, conn)
	if e != nil {
		panic(e)
	}
	return db
}

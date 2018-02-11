package model

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"regexp"

	"github.com/godcong/oauth2server/base"
	"golang.org/x/crypto/bcrypt"
)

type AccountType int

const (
	ACCOUNT_TYPE_NONE   = 0
	ACCOUNT_TYPE_EMAIL  = 1 << iota
	ACCOUNT_TYPE_MOBILE = ACCOUNT_TYPE_EMAIL << iota
	ACCOUNT_TYPE_UNAME  = ACCOUNT_TYPE_EMAIL << iota
)

type RegisterType string

type User struct {
	BaseModel
	Username      string
	Nickname      string
	Mobile        string
	Mail          string
	Password      string
	Salt          string
	Status        int         //正常 锁定 ...
	AccountType   AccountType //账户类型
	RegisterType  string      //定义注册来源
	ErrorTimes    int
	LoginIp       string //登录地址
	LoginAt       string //登录设备
	LoginFrom     string //登录客户
	LastLoginIp   string
	LastLoginAt   string
	LastLoginFrom string
}

func init() {
	AddModel("User", User{})
}

//create a user with salt
func NewUser() *User {
	var user User
	user.Salt = string(base.GenerateRandomString(16, base.T_RAND_ALL))
	return &user
}

func (u *User) GenerateMD5Password(s string) {
	m5 := md5.New()
	m5.Write([]byte(s))
	m5.Write([]byte(u.Salt))
	st := m5.Sum(nil)
	u.Password = hex.EncodeToString(st)
}

func (u *User) GenerateBCryptPassword(s string) {
	if p, e := bcrypt.GenerateFromPassword([]byte(s+u.Salt), 10); e == nil {
		u.Password = string(p)
	}
}

func (u *User) GenerateMD5BCryptPassword(s string) {
	u.GenerateMD5Password(s + u.Salt)
	u.GenerateBCryptPassword(u.Password + u.Salt)
}

func (u *User) FirstByID(id string) {
	FirstById(u, id)
	fmt.Println(*u)
}

func (u *User) FirstWhere(query string, args ...interface{}) {
	FirstWhere(u, query, args)
	fmt.Println(*u)
}

func VerifyAccountType(uname string) AccountType {

	if b, e := VerifyUsername(uname); b == true && e == nil {
		return ACCOUNT_TYPE_UNAME
	}

	if b, e := VerifyMail(uname); b == true && e == nil {
		return ACCOUNT_TYPE_EMAIL
	}

	if b, e := VerifyMobile(uname); b == true && e == nil {
		return ACCOUNT_TYPE_MOBILE
	}

	return ACCOUNT_TYPE_NONE
}

func VerifyUsername(s string) (bool, error) {
	//b, e := regexp.Match("^[A-Za-z]{1}.\\w+$", []byte(s))
	return regexp.Match("^[A-Za-z]{1}.\\w+$", []byte(s))
}

func VerifyMail(s string) (bool, error) {
	//b, e := regexp.Match("@.*", []byte(s))
	return regexp.Match("@.*", []byte(s))

}

func VerifyMobile(s string) (bool, error) {
	//b, e := regexp.Match("^\\d{11}$", []byte(s))
	return regexp.Match("^\\d{11}$", []byte(s))

}

func (u *User) GetUser(s string, t AccountType) (ret bool) {
	ret = true
	if t == ACCOUNT_TYPE_MOBILE {
		u.GetUserByMobile(s)
	} else if t == ACCOUNT_TYPE_EMAIL {
		u.GetUserByMail(s)
	} else if t == ACCOUNT_TYPE_UNAME {
		u.GetUserByUsername(s)
	} else {
		ret = false
	}
	return
}

func (u *User) GetUserByUsername(s string) {
	Gorm().Where("username = ?", s).First(u)
}

func (u *User) GetUserByMobile(s string) {
	Gorm().Where("mobile = ?", s).First(u)
}

func (u *User) GetUserByMail(s string) {
	Gorm().Where("mail = ?", s).First(u)
}

//更新登录信息
func (u *User) UpdateSignInfo(ip, at, from string) {
	u.LastLoginAt = u.LoginAt
	u.LastLoginFrom = u.LoginFrom
	u.LastLoginIp = u.LoginIp
	u.LoginAt = at
	u.LoginFrom = from
	u.LoginIp = ip
}

//密码验证
func (u *User) VerifyPassword(s string) bool {
	if u.ErrorTimes > 5 {
		return false
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(s+u.Salt)); err != nil {
		u.ErrorTimes += 1
		return false
	}
	u.ErrorTimes = 0
	return true
}

//func (u *User) Save() *gorm.DB {
//	return Gorm().Save(u)
//}

func (u *User) IsNull() bool {
	return IsNull(u.ID)
}

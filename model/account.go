package model

import (
	"fmt"

	"github.com/satori/go.uuid"
)

type Account struct {
	BaseModel
	AccountName     string
	AccountPassword string
	AccountId       uuid.UUID `gorm:"type:varchar(36)"`
	UserID          uuid.UUID `gorm:"type:varchar(36)"`
}

func init() {
	AddModel("Account", Account{})
}

func (u *Account) FirstByID(id string) {
	FirstById(u, id)
	fmt.Println(*u)
}

func (u *Account) FirstWhere(query string, args ...interface{}) {
	FirstWhere(u, query, args)
	fmt.Println(*u)
}

func (u *Account) GetAccountByName(s string) {
	Gorm().Where("username = ?", s).First(u)
}

package model

import (
	"github.com/godcong/oauth2server/base"

	"github.com/satori/go.uuid"
)

type Authorize struct {
	BaseModel
	SubID        string
	UserID       uuid.UUID `gorm:"type:varchar(36)"`
	User         User      `gorm:"ForeignKey:UserID"`
	ClientID     uuid.UUID `gorm:"type:varchar(36)"`
	Client       Client    `gorm:"ForeignKey:ClientID"`
	AccessToken  string
	RefreshToken string
	TokenType    string
	ExpiresIn    int
}

//noinspection ALL
const PREFIX_OPENID = "mn_"

func init() {
	AddModel("Authorize", Authorize{})
}

//0 cid,1 uid,2 prefix
func GenerateOpenID(args ...string) (oid string) {
	if len(args) > 2 {
		oid = args[2] + base.GenerateMD5(args[0], args[1])
		return
	}
	oid = PREFIX_OPENID + base.GenerateMD5(args[0], args[1])
	return

}

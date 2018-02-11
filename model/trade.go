package model

type Trade struct {
	BaseModel
}

func init() {
	AddModel("Trade", Trade{})
}

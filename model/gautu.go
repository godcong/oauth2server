package model

//something about the system data
type Gautu struct {
	BaseModel
	//gorm.Model
}

func init() {
	AddModel("Gautu", Client{})
}

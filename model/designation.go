package model

type Designation struct {
	Desigcode int    `json:"desigcode" gorm:"type:int;not null;primary_key;auto_increment;unique"`
	Designame string `json:"designame" gorm:"type:varchar(50);not null;unique"`
	Sdesig    string `json:"sdesig" gorm:"type:varchar(10);not null;unique"`
}

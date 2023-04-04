package model

type Acchead struct {
	Id         int    `json:"id"`
	Accid      int    `gorm:"primaryKey;autoIncrement:false" json:"accid"`
	Acccode    string `json:"acccode"`
	Parent     int    `json:"parent"`
	Name       string `json:"name"`
	Lr         string `json:"lr"`
	Category   string `json:"category"`
	Createdate string `json:"createdate"`
	Topparent  int    `json:"topparent"`
	Depth      int    `json:"depth"`
}

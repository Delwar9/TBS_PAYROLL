package model

type Salarystructure struct {
	Id               int     `json:"id"`
	Empcode          int     `json:"empcode"`
	Refno            int     `json:"refno"`
	Incrementno      int     `json:"incrementno"`
	Incrementamount  float64 `json:"incrementamount"`
	Desigcode        int     `json:"desigcode"`
	Branchcode       int     `json:"branchcode"`
	Deptcode         int     `json:"deptcode"`
	Bankid           int     `json:"bankid"`
	Pfbankid         int     `json:"pfbankid"`
	Consolited       float64 `json:"consolited"`
	Basic            float64 `json:"basic" sql:"type:float(10,2)"`
	Houserent        float64 `json:"houserent"`
	Conveyance       float64 `json:"conveyance"`
	Medical          float64 `json:"medical"`
	Entertainment    float64 `json:"entertainment"`
	Housemaint       float64 `json:"housemaint"`
	Incometax        float64 `json:"incometax"`
	Bonusrate        float64 `json:"bonusrate"`
	Arrear           float64 `json:"arrear"`
	Cpf              float64 `json:"cpf"`
	Groupins         float64 `json:"groupins"`
	Cpfloan          float64 `json:"cpfloan"`
	Stamp            float64 `json:"stamp"`
	Pfund            float64 `json:"pfund"`
	Sal_scale        string  `json:"sal_scale"`
	Seniority_serial int     `json:"seniority_serial"`
	Telephone        string  `json:"telephone"`
	Incentive        float64 `json:"incentive"`
	Specialallow     float64 `json:"specialallow"`
	Overtime         float64 `json:"overtime"`
	Food             float64 `json:"food"`
	Salaryadv        float64 `json:"salaryadv"`
	Otherallow       float64 `json:"otherallow"` // added
	Otheradv         float64 `json:"otheradv"`
	Carallow         float64 `json:"carallow"`
	Specialallow1    float64 `json:"specialallow1"` // exclude
	Binder_wf        float64 `json:"binder_wf"`     // modify
	Leaserent        float64 `json:"leaserent"`
	Dareness         float64 `json:"dareness"`    // added
	Specialda        float64 `json:"specialda"`   // added
	Extraallow       float64 `json:"extraallow"`  // added
	Technical        float64 `json:"technical"`   // added
	Mobile           float64 `json:"mobile"`      // added
	Pubsalary        float64 `json:"pubsalary"`   // added
	Business         float64 `json:"business"`    // added
	Charge           float64 `json:"charge"`      // added
	Eyeallow         float64 `json:"eyeallow"`    // added
	Cosecretary      float64 `json:"cosecretary"` // added
	Grosssalary      float64 `json:"grosssalary"` // added
}

type Salarystructure_archive struct {
	Salarystructure

	Changedate   string `json:"changedate"`
	Changeuserid string `json:"changeuserid"`
	Flag_ed_del  string `json:"flag_ed_del"`
	Trackid      int    `json:"trackid"`
}

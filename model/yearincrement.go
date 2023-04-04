package model

type Yearincrement struct {
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
	Basic            float64 `json:"basic"`
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
	Otherallow       float64 `json:"otherallow"`
	Otheradv         float64 `json:"otheradv"`
	Carallow         float64 `json:"carallow"`
	Specialallow1    float64 `json:"specialallow1"`
	Binder_wf        float64 `json:"binder_wf"`
	Leaserent        float64 `json:"leaserent"`
	Dareness         float64 `json:"dareness"`
	Specialda        float64 `json:"specialda"`
	Extraallow       float64 `json:"extraallow"`
	Technical        float64 `json:"technical"`
	Mobile           float64 `json:"mobile"`
	Pubsalary        float64 `json:"pubsalary"`
	Business         float64 `json:"business"`
	Charge           float64 `json:"charge"`
	Eyeallow         float64 `json:"eyeallow"`
	Cosecretary      float64 `json:"cosecretary"`
	Grosssalary      float64 `json:"grosssalary"`
	Month            int     `json:"month"`           // added
	Year             int     `json:"year"`            // added
	Joindate         string  `json:"joindate"`        // added
	Probationperiod  int     `json:"probationperiod"` // added
	Confirmdate      string  `json:"confirmdate"`     // added
	Entry_user       string  `json:"entry_user"`      // added
	Accno            string  `json:"accno"`           // added
	Active_salary    int     `json:"active_salary"`   // added
}

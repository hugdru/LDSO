package data

type Main_Group struct {
	Id     int    `json:"_id" bson:"_id,omitempty"`
	Name   string `json:"name"`
	Weight int    `json:"weight"`
}

type Sub_Group struct {
	Id         int    `json:"_id" bson:"_id,omitempty"`
	Name       string `json:"name"`
	Weight     int    `json:"weight"`
	Main_Group int    `json:"main_group"`
}

type Criterion struct {
	Id          int    `json:"_id" bson:"_id,omitempty"`
	Name        string `json:"name"`
	Weight      int    `json:"weight"`
	Legislation string `json:"legislation"`
	Sub_Group   int    `json:"sub_group"`
}

type Accessibility struct {
	Id        int    `json:"_id" bson:"_id,omitempty"`
	Name      string `json:"name"`
	Weight    int    `json:"weight"`
	Criterion int    `json:"criterion"`
}

type Owner struct {
	Name string `json:"name"`
}

type Property struct {
	Id   int    `json:"_id" bson:"_id,omitempty"`
	Name string `json:"name"`
}

type Audit struct {
	Id       int              `json:"_id" bson:"_id,omitempty"`
	Property int              `json:"property"`
	Rating   int              `json:"rating"`
	Criteria []AuditCriterion `json:"criteria"`
}

type AuditCriterion struct {
	Criterion int `json:"criterion"`
	Value     int `json:"value"`
}

type Note struct {
	Data       []string  `json:"data"`
	Criterion  Criterion `json:"criterion"`
	Image_path []string  `json:"image_path"`
}

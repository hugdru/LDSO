package data

type Main_Group struct {
	Id int`json:"_id" bson:"_id,omitempty"`
	Name string `json:"name"`
	Weight int `json:"weight"`
}

type Sub_Group struct {
	Id int `json:"_id" bson:"_id,omitempty"`
	Name string `json:"name"`
	Weight int `json:"weight"`
	Main_Group int `json:"main_group"`
}

type Criterion struct {
	Id int `json:"_id" bson:"_id,omitempty"`
	Name string `json:"name"`
	Weight int `json:"weight"`
	Legislation string `json:"legislation"`
	Sub_Group int `json:"sub_group"`
}

type Accessibility struct {
	Id int `json:"_id" bson:"_id,omitempty"`
	Name string `json:"name"`
	Weight int `json:"weight"`
	Criterion int `json:"criterion"`
}

type Owner struct {
	Name string `json:"name"`
}

type Property struct {
	Name string `json:"name"`
	Owner Owner `json:"owner"`
	Image_path string `json:"image_path"`
}

type Note struct {
	Data []string `json:"data"`
	Criterion Criterion `json:"criterion"`
	Image_path []string `json:"image_path"`
}

type Audit struct {
	Property Property `json:"property"`
	Notes []Note `json:"notes"`
	Value int `json:"value"`
}

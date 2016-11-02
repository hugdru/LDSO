package data

type Main_Group struct {
	Id uint `json:"_id" bson:"_id,omitempty"`
	Name string `json:"name"`
	Weight uint `json:"weight"`
}

type Sub_Group struct {
	Id uint `json:"_id" bson:"_id,omitempty"`
	Name string `json:"name"`
	Weight uint `json:"weight"`
	Main_Group uint `json:"main_group"`
}

type Criterion struct {
	Id uint `json:"_id" bson:"_id,omitempty"`
	Name string `json:"name"`
	Weight uint `json:"weight"`
	Legislation string `json:"legislation"`
	Sub_Group uint `json:"sub_group"`
}

type Accessibility struct {
	Id uint `json:"_id" bson:"_id,omitempty"`
	Name string `json:"name"`
	Weight uint `json:"weight"`
	Criterion uint `json:"criterion"`
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
	Value uint `json:"value"`
}

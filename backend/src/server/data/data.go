package data

type Group struct {
	Name string `json:name`
	Weight uint `json:weight`
}

type Sub_Group struct {
	Name string `json:name`
	Weight uint `json:weight`
}

type Accessibility string

type Criterion struct {
	Name string `json:name`
	Accessibility Accessibility `json:accessibility`
	Legislation bool `json:legislation`
}

type Group_Set struct {
	Group Group `json:group`
	Sub_Groups []Sub_Group_Set `json:sub_groups`
}

type Sub_Group_Set struct {
	Sub_Group Sub_Group `json:sub_group`
	Criteria []Criterion `json:criteria`
}

type Owner struct {
	Name string `json:name`
}

type Property struct {
	Name string `json:"name"`
	Owner Owner `json:"owner"`
	Image_path string `json:"image_path"`
}

type Note struct {
	Data string `json:data`
	Criterion Criterion `json:sub_criterion`
	Image_path string `json:image_path`
}

type Evaluation struct {
	Property Property `json:property`
	Notes []Note `json:notes`
	Value uint `json:value`
}

func (this *Group_Set) SetSub(subs ...Sub_Group) {
	this.Sub_Group = nil
	for _, sub := range subs {
		this.Sub_Group = append(this.Sub_Group, sub)
	}
}

func (this *Group_Set) AppendSub(sub Sub_Group) {
	this.Sub_Group = append(this.Sub_Group, sub)
}

func (this *Sub_Group_Set) SetSub(subs ...Sub_Group) {
	this.Criteria = nil
	for _, sub := range subs {
		this.Criteria = append(this.Criteria, sub)
	}
}

func (this *Sub_Group_Set) AppendSub(sub Sub_Group) {
	this.Criteria = append(this.Criteria, sub)
}

package data

type Group struct {
	Name string `json:"name"`
	Weight uint `json:"weight"`
	Sub_Groups []Sub_Group `json:"sub_groups"`
}

type Sub_Group struct {
	Name string `json:"name"`
	Weight uint `json:"weight"`
	Criteria []Criterion `json:"criteria"`
}

type Criterion struct {
	Name string `json:"name"`
	Weight uint `json:"weight"`
	Accessibilities []Accessibility `json:"accessibilities"`
	Legislation string `json:"legislation"`
}

type Accessibility struct {
	Name string `json:"name"`
	Weight uint `json:"weight"`
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
	Data string `json:"data"`
	Criterion Criterion `json:"sub_criterion"`
	Image_path string `json:"image_path"`
}

type Evaluation struct {
	Property Property `json:"property"`
	Notes []Note `json:"notes"`
	Value uint `json:"value"`
}

func (this *Sub_Group) SetCriteria(subs ...Criterion) {
	this.Criteria = nil
	for _, sub := range subs {
		this.Criteria = append(this.Criteria, sub)
	}
}

func (this *Sub_Group) AppendCriteria(sub Criterion) {
	this.Criteria = append(this.Criteria, sub)
}

func (this *Group) SetSubs(subs ...Sub_Group) {
	this.Sub_Groups = nil
	for _, sub := range subs {
		this.Sub_Groups = append(this.Sub_Groups, sub)
	}
}

func (this *Group) AppendSub(sub Sub_Group) {
	this.Sub_Groups = append(this.Sub_Groups, sub)
}

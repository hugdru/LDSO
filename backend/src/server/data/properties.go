ackage data

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
	Sub_Criterion Sub_Criterion `json:sub_criterion`
	Image_path string `json:image_path`
}

type Evaluation struct {
	Property Property `json:property`
	Criteria_Set []Criterion_Set `json:criteria_set`
	Notes []Note `json:notes`
	Value uint `json:value`
}

func (this *Evaluation) SetCriteria(crit ...Criterion_Set) {
	this.Criteria_Set = nil
	for _, sub := range crit {
		this.Criteria_Set = append(this.Criteria_Set, sub)
	}
}

func (this *Evaluation) AddCriteriaSet(sub Criterion_Set) {
	this.Criteria_Set = append(this.Criteria_Set, sub)
}

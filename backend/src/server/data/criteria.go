package data

type Criterion struct {
	Name string `json:"name"`
	Weight uint `json:"weight"`
}

type Sub_Criterion struct {
	Name string `json:name`
	Weight uint `json:weight`
}

type Criterion_Set struct {
	Criterion Criterion `json:criterion`
	Sub_Criteria []Sub_Criterion `json:sub_criteria`
}

func (this *Criterion_Set) SetSub_Criterion(subs ...Sub_Criterion) {
	this.Sub_Criteria = nil
	for _, sub := range subs {
		this.Sub_Criteria = append(this.Sub_Criteria, sub)
	}
}

func (this *Criterion_Set) AppendSub_Criterion(sub Sub_Criterion) {
	this.Sub_Criteria = append(this.Sub_Criteria, sub)
}

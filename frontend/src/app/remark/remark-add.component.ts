import {Component, Input, Output, OnInit, EventEmitter } from "@angular/core";

import {Criterion} from "criterion/criterion";
import {Remark} from "remark/remark";

@Component({
    selector: 'remark-add',
    templateUrl: 'html/remark-add.component.html',
    styleUrls: ['../main-group/main-group.component.css'],
    providers: []
})

export class RemarkAddComponent {
	remark: Remark;

	@Input() criterion: Criterion;
	@Output() add = new EventEmitter<Remark>();

	constructor() { }

	ngOnInit(): void {
		this.remark = new Remark();
	}

	onAdd(remark: Remark): void {
		this.remark.idCriterion = this.criterion._id;
		this.add.emit(remark);
	}
}

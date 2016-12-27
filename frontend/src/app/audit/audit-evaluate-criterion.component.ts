import {Component, OnInit, Input} from "@angular/core";
import {MdCheckbox} from '@angular/material';

import {Criterion} from "criterion/criterion";
import {Remark} from "remark/remark";

@Component({
    selector: 'audit-evaluate-criterion',
    templateUrl: 'html/audit-evaluate-criterion.component.html',
    styleUrls: ['../main-group/main-group.component.css']
    // providers: [MainGroupService, SubGroupService, CriterionService]
})

export class AuditEvaluateCriterionComponent {
	@Input() criteria: Criterion[];

	uncheckedCriteria: Criterion[] = [];

	remark: Remark;
	remarks: Remark[];
	selectedAdd: boolean = false;
	a: boolean = false;

    ngOnInit(): void {
		this.remarks = [];
    }

	checkedNoCriterion(criterion: Criterion): void {
		console.log("hello");
		// this.uncheckedCriteria.push(criterion);
	}

	uncheckedNoCriterion(criterion: Criterion, change: boolean): void {
		console.log("hello " + criterion.name + "  " + change + " " + this.a);
	}

	checkCriterion(criterion: Criterion): boolean {
		return this.uncheckedCriteria.includes(criterion);
	}

	selectAdd(): void {
		this.selectedAdd = true;
	}

	onAdd(remark: Remark): void {
		if (remark != null) {
			this.remarks.push(remark);
		}
		this.selectedAdd = false;
	}
}

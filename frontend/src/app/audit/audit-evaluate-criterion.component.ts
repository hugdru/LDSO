import {Component, OnInit, Input} from "@angular/core";

import {Criterion} from "criterion/criterion";
import {Remark} from "remark/remark";

@Component({
    selector: 'audit-evaluate-criterion',
    templateUrl: 'html/audit-evaluate-criterion.component.html',
    styleUrls: ['../main-group/main-group.component.css'],
    // providers: [MainGroupService, SubGroupService, CriterionService]
})

export class AuditEvaluateCriterionComponent {
	@Input() criteria: Criterion[];

	remark: Remark;
	remarks: Remark[];
	selectedAdd: boolean = false;

    ngOnInit(): void {
		this.remarks = [];
    }

	checkedNoCriterion(criterion: Criterion): void {

	}

	uncheckedNoCriterion(criterion: Criterion): void {

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

import {Component, OnInit, Input} from "@angular/core";

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

	unselectedCriteria: Criterion[] = [];

	remark: Remark;
	remarks: Remark[];
	selectedAdd: boolean = false;
	checked: boolean = false;
	b: number = 1;

    ngOnInit(): void {
		this.remarks = [];
    }

	changedCheckbox(): void {
		this.checked = !this.checked;
	}

	setCheckboxValue(criterion: Criterion): void {
		this.checked = this.checkUnselected(criterion);
	}

	submitCriterion(criterion: Criterion): void {
		if (this.checked && !this.checkUnselected(criterion)) {
			this.unselectedCriteria.push(criterion);
		} else if (!this.checked && this.checkUnselected(criterion)) {
			let index = this.unselectedCriteria.indexOf(criterion, 0);
			if (index > -1) {
				this.unselectedCriteria.splice(index, 1);
			}
		}
	}

	checkUnselected(criterion: Criterion): boolean {
		return this.unselectedCriteria.includes(criterion);
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

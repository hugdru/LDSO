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
	idAudit = 1;
	id: number = 0;
	selectedId; number = -1;
	uncheckedCriteria: Criterion[] = [];
	remarks: Remark[];
	selectedAdd: boolean = false;

	@Input() criteria: Criterion[];

    ngOnInit(): void {
		this.remarks = [];
    }

	checkedNoCriterion(criterion: Criterion): void {
		// console.log("hello");
			// this.uncheckedCriteria.push(criterion);
		}

	uncheckedNoCriterion(criterion: Criterion): void {

	}

	checkCriterion(criterion: Criterion): boolean {
		return this.uncheckedCriteria.includes(criterion);
	}

	selectAdd(): void {
		this.selectedAdd = true;
	}

	onAdd(remark: Remark): void {
		if (remark != null) {
			remark.idAudit = this.idAudit;
			remark.id = this.id++;
			this.remarks.push(remark);
		}
		this.selectedAdd = false;
	}

	selectedRemark(id: number): void {
		this.selectedId = id;
	}
}

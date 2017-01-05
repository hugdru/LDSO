import {Component, OnInit, Input} from "@angular/core";

import {Criterion} from "criterion/criterion";
import {AuditCriterion} from "audit/audit";
import {Remark} from "remark/remark";
import {AuditService} from 'audit/service/audit.service';

@Component({
    selector: 'audit-evaluate-criterion',
    templateUrl: 'html/audit-evaluate-criterion.component.html',
    styleUrls: ['../main-group/main-group.component.css'],
    providers: [AuditService]
})

export class AuditEvaluateCriterionComponent {
	@Input() criteria: Criterion[];
	@Input() auditId: number;

	id: number = 0;
	selectedId; number = -1;
	unselectedCriteria: Criterion[] = [];
	remarks: Remark[];
	selectedAdd: boolean = false;
	checked: boolean = false;
	rating: number;

	constructor(private auditService: AuditService) {

	}

    ngOnInit(): void {
		this.remarks = [];
		console.log("auditid: " + this.auditId);
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
			let auditCriterion: AuditCriterion;
			auditCriterion.criterion = criterion.id;
			auditCriterion.rating = this.rating;
			// this.auditService.setAuditCriterion(auditCriterion).subscribe();
		}
		
	}

	checkUnselected(criterion: Criterion): boolean {
		return this.unselectedCriteria.includes(criterion);
	}

	checkCriterion(criterion: Criterion): boolean {
		return this.unselectedCriteria.includes(criterion);
	}

	selectAdd(): void {
		this.selectedAdd = true;
	}

	onAdd(remark: Remark): void {
		if (remark != null) {
			remark.idAudit = this.auditId;
			remark.id = this.id++;
			this.remarks.push(remark);
		}
		this.selectedAdd = false;
	}

	selectedRemark(id: number): void {
		this.selectedId = id;
	}
}

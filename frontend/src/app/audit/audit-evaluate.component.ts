import {Component, OnInit, Input} from "@angular/core";

import {MainGroupService} from "main-group/service/main-group.service";
import {MainGroup} from "main-group/main-group";
import {SubGroupService} from "sub-group/service/sub-group.service";
import {SubGroup} from "sub-group/sub-group";
import {CriterionService} from "criterion/service/criterion.service";
import {Criterion} from "criterion/criterion";

@Component({
    selector: 'audit-evaluate',
    templateUrl: 'html/audit-evaluate.component.html',
    styleUrls: ['../main-group/main-group.component.css'],
    providers: [MainGroupService, SubGroupService, CriterionService]
})

export class AuditEvaluateComponent implements OnInit {
    @Input() selectedSubGroups: SubGroup[];
    @Input() auditId: number;

    mainGroups: MainGroup[] = [];
    mainGroupsId: number[] = [];
	subGroups: SubGroup[];
    criteria: Criterion[];

    constructor(private mainGroupService: MainGroupService,
			private criterionService: CriterionService) {
    }

    ngOnInit(): void {
        this.findMainGroups();
        this.getMainGroups();
    }

    getMainGroups(): void {
        for (let mainGroupId of this.mainGroupsId) {
            this.mainGroupService.getMainGroup(mainGroupId)
                    .subscribe(data => this.mainGroups.push(data));
        }
    }

    findMainGroups(): void {
        for (let subGroup of this.selectedSubGroups) {
            if (!this.mainGroupsId.includes(subGroup.idMaingroup)) {
                this.mainGroupsId.push(subGroup.idMaingroup);
            }
        }
    }

    selected(object: Object): void {
        if ((<SubGroup>object).idMaingroup !== undefined) {
			this.showCriteria(<SubGroup>object);
        }
        else {
			this.initSubGroups(<MainGroup>object);
			this.criteria = [];
        }
    }

    showCriteria(subGroup: SubGroup): void {
        this.initCriteria(subGroup);
    }

    initCriteria(subGroup: SubGroup): void {
        this.criterionService.getSomeCriteria("idSubgroup",
                subGroup.id).subscribe(data => this.criteria = data);
    }

	initSubGroups(mainGroup: MainGroup): void {
        this.subGroups = [];
        for(let subGroup of this.selectedSubGroups) {
            if (subGroup.idMaingroup == mainGroup.id) {
                this.subGroups.push(subGroup);
            }
        }
	}

	finishedAudit(){
		//TODO
	}
}

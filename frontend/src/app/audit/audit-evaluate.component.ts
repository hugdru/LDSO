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

    mainGroups: MainGroup[] = [];
    mainGroupsId: number[] = [];
	subGroups: SubGroup[];
    criteria: Criterion[];

    constructor(private mainGroupService: MainGroupService,
			private subGroupService: SubGroupService,
			private criterionService: CriterionService) {
    }

    ngOnInit(): void {
        this.findMainGroups();
        this.getMainGroups();
    }

    getMainGroups(): void {
        for (let mainGroupId of this.mainGroupsId) {
            this.mainGroupService.getMainGroup("_id", "int", mainGroupId)
                    .subscribe(data => this.mainGroups.push(data));
        }
    }

    findMainGroups(): void {
        for (let subGroup of this.selectedSubGroups) {
            if (!this.mainGroupsId.includes(subGroup.main_group)) {
                this.mainGroupsId.push(subGroup.main_group);
            }
        }
    }

    selected(object: Object): void {
        if ((<SubGroup>object).main_group !== undefined) {
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
        this.criterionService.getSomeCriteria("sub_group", "int",
                subGroup._id).subscribe(data => this.criteria = data);
    }

	initSubGroups(mainGroup: MainGroup): void {
        this.subGroupService.getSomeSubGroups("main_group", "int",
                mainGroup._id).subscribe(data => this.subGroups = data);
	}

}

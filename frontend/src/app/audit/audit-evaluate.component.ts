import {Component, OnInit, Input} from "@angular/core";
import {MainGroupService} from "main-group/service/main-group.service";
import {MainGroup} from "main-group/main-group";
import {SubGroup} from "sub-group/sub-group";
import {CriterionService} from "../criterion/service/criterion.service";
import {Criterion} from "../criterion/criterion";

@Component({
    selector: 'audit-evaluate',
    templateUrl: 'html/audit-evaluate.component.html',
    styleUrls: ['../main-group/main-group.component.css'],
    providers: [MainGroupService, CriterionService]
})

export class AuditEvaluateComponent implements OnInit {
    @Input() selectedSubGroups: SubGroup[];
    mainGroups: MainGroup[] = [];
    mainGroupsId: number[] = [];
    selectedMainGroup: MainGroup;
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

    showSubGroups(mainGroup: MainGroup): void {
        this.selectedMainGroup = mainGroup;
        this.criteria = [];
    }

    showCriteria(subGroup: SubGroup): void {
        this.initCriteria(subGroup);
    }

    initCriteria(subGroup: SubGroup): void {
        this.criterionService.getSomeCriteria("sub_group", "int",
                subGroup._id).subscribe(data => this.criteria = data);
    }

	checkedNoCriterion(criterion: Criterion): void {

	}

	uncheckedNoCriterion(criterion: Criterion): void {

	}
}

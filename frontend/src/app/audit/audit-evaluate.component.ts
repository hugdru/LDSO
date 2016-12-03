import {Component, OnInit, Input} from "@angular/core";
import {MainGroupService} from "main-group/service/main-group.service";
import {MainGroup} from "main-group/main-group";
import {SubGroupService} from "sub-group/service/sub-group.service";
import {SubGroup} from "sub-group/sub-group";
import {CriterionService} from "criterion/service/criterion.service";
import {Criterion} from "criterion/criterion";
import {Remark} from "remark/remark";

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
	selectedAdd: boolean = false;
	remark: Remark;
	remarks: Remark[];

    constructor(private mainGroupService: MainGroupService,
			private subGroupService: SubGroupService,
			private criterionService: CriterionService) {
    }

    ngOnInit(): void {
        this.findMainGroups();
        this.getMainGroups();
		this.remarks = [];
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
        if ((<SubGroup>object).main_group !== undefined) {
			this.showCriteria(<SubGroup>object);
        }
        // else if ((<Criterion>object).sub_group !== undefined) {

        // }
        else {
			this.initSubGroups(<MainGroup>object);
			this.criteria = [];
        }
    }

    showCriteria(subGroup: SubGroup): void {
        this.initCriteria(subGroup);
    }

    initCriteria(subGroup: SubGroup): void {
        this.criterionService.getSomeCriteria("idSubGroup",
                subGroup.id).subscribe(data => this.criteria = data);
    }

	initSubGroups(mainGroup: MainGroup): void {
        this.subGroupService.getSomeSubGroups("main_group", "int",
                mainGroup._id).subscribe(data => this.subGroups = data);
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

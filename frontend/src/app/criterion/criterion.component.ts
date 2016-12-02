import {
        Component,
        OnInit,
        OnChanges,
        SimpleChanges,
        Input
} from "@angular/core";
import {CriterionService} from "criterion/service/criterion.service";
import {Criterion} from "criterion/criterion";
import {SubGroup} from "sub-group/sub-group";

@Component({
    selector: 'criterion',
    templateUrl: 'html/criterion.component.html',
    styleUrls: ['../main-group/main-group.component.css'],
    providers: [CriterionService]
})

export class CriterionComponent implements OnInit, OnChanges {
    criteria: Criterion[];
    parentCriterion: Criterion;

    @Input() parentSubGroup: SubGroup;

    constructor(private criterionService: CriterionService) {
    }

    ngOnChanges(changes: SimpleChanges): void {
        for (let i in changes) {
            this.initCriteria(changes[i].currentValue.id);
            this.parentCriterion = undefined;
        }
    }

    ngOnInit() {
        this.initCriteria(this.parentSubGroup.id);
    }

    initCriteria(subGroupId: number): void {
        this.criterionService
                .getSomeCriteria("idSubgroup", subGroupId)
                .subscribe(data => this.criteria = data);

    }

    onDelete(criterion: Criterion): void {
        this.criterionService.removeCriterion(criterion.id).subscribe();
    }

    onShow(criterion: Criterion): void {
        this.parentCriterion = criterion;
    }

}

import {Component, Input, Output, EventEmitter, OnInit} from "@angular/core";
import {CriterionService} from "criterion/service/criterion.service";
import {Criterion} from "criterion/criterion";

@Component({
    selector: 'criterion-edit',
    templateUrl: '../audit-template/html/audit-template-edit.component.html',
    styleUrls: ['../audit-template/audit-template-edit.component.css'],
    providers: [CriterionService]
})

export class CriterionEditComponent implements OnInit {
    backupCriterion: Criterion;
    goodPractice: Boolean;

    @Input() objType: string;
    @Input() selectedObject: Criterion;
    @Input() weight: number;
    @Output() onAction = new EventEmitter();

    constructor(private criterionService: CriterionService) {
    }

    ngOnInit(): void {
        this.backupCriterion = new Criterion();
        this.backupCriterion.name = this.selectedObject.name;
        this.backupCriterion.weight = this.selectedObject.weight;
        this.goodPractice = this.selectedObject.legislation == '';
    }

    pressed(updatedCriterion: Criterion): void {
        if (updatedCriterion) {
            this.updateCriterion();
        } else {
            this.selectedObject.name = this.backupCriterion.name;
            this.selectedObject.weight = this.backupCriterion.weight;
        }
        this.onAction.emit();
    }

    updateCriterion(): void {
        this.criterionService.updateCriterion(this.selectedObject)
                .subscribe();
    }

    checkPercentage(): boolean {
        return this.selectedObject.weight + this.weight != 100;
    }

}

import {Component, Input, Output, EventEmitter, OnInit} from "@angular/core";
import {CriterionService} from "criterion/service/criterion.service";
import {Criterion} from "criterion/criterion";
import {SubGroup} from "sub-group/sub-group";

@Component({
    selector: 'criterion-add',
    templateUrl: '../main-group/html/main-group-edit.component.html',
    styleUrls: ['../main-group/main-group-edit.component.css'],
    providers: [CriterionService]
})

export class CriterionAddComponent implements OnInit {
    selectedObject: Criterion;
    goodPratice: boolean = false;

    @Input() objType: string;
    @Input() subGroup: SubGroup;
    @Input() weight: number;
    @Output() onAdd = new EventEmitter<Criterion>();

    constructor(private criterionService: CriterionService) {

    }

    ngOnInit(): void {
        this.selectedObject = new Criterion();
    }

    pressed(newCriterion: Criterion): void {
        if (newCriterion) {
            this.addCriterion();
        }
        this.onAdd.emit(newCriterion);
    }

    addCriterion(): void {
        this.selectedObject.idSubgroup = this.subGroup.id;
        this.criterionService.setCriterion(this.selectedObject).subscribe(
                response => this.selectedObject.id = response.json()
        );
    }

    checkPercentage(): boolean {
        return this.selectedObject.weight + this.weight > 100;
    }
}

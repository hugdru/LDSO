import {Component, Input, Output, EventEmitter, OnInit} from "@angular/core";
import {SubGroupService} from "sub-group/service/sub-group.service";
import {SubGroup} from "sub-group/sub-group";

@Component({
    selector: 'sub-group-edit',
    templateUrl: '../audit-template/html/audit-template-edit.component.html',
    styleUrls: ['../audit-template/audit-template-edit.component.css'],
    providers: [SubGroupService]
})

export class SubGroupEditComponent implements OnInit {
    backupSubGroup: SubGroup;

    @Input() selectedObject: SubGroup;
    @Input() weight: number;
    @Output() onAction = new EventEmitter();

    constructor(private subGroupService: SubGroupService) {

    }

    ngOnInit(): void {
        this.backupSubGroup = new SubGroup();
        this.backupSubGroup.name = this.selectedObject.name;
        this.backupSubGroup.weight = this.selectedObject.weight;
    }

    pressed(updatedSubGroup: SubGroup): void {
        if (updatedSubGroup) {
            this.updateSubGroup();
        } else {
            this.selectedObject.name = this.backupSubGroup.name;
            this.selectedObject.weight = this.backupSubGroup.weight;
        }
        this.onAction.emit();
    }

    updateSubGroup(): void {
        this.subGroupService.updateSubGroup(this.selectedObject)
                .subscribe();
    }

    checkPercentage(): boolean {
        return this.selectedObject.weight + this.weight != 100;
    }

}

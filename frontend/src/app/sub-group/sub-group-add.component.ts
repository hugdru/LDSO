import {Component, Input, Output, EventEmitter, OnInit} from "@angular/core";
import {SubGroupService} from "sub-group/service/sub-group.service";
import {SubGroup} from "sub-group/sub-group";
import {MainGroup} from "main-group/main-group";

@Component({
    selector: 'sub-group-add',
    templateUrl: '../audit-template/html/audit-template-edit.component.html',
    styleUrls: ['../audit-template/audit-template-edit.component.css'],
    providers: [SubGroupService]
})

export class SubGroupAddComponent implements OnInit {
    selectedObject: SubGroup;

    @Input() mainGroup: MainGroup;
    @Input() weight: number;
    @Output() onAdd = new EventEmitter<SubGroup>();

    constructor(private subGroupService: SubGroupService) {

    }

    ngOnInit(): void {
        this.selectedObject = new SubGroup();
    }

    pressed(newSubGroup: SubGroup): void {
        if (newSubGroup) {
            this.addSubGroup();
        }
        this.onAdd.emit(newSubGroup);
    }

    addSubGroup(): void {
        this.selectedObject.idMaingroup = this.mainGroup.id;
        this.subGroupService.setSubGroup(this.selectedObject).subscribe(
                response => this.selectedObject.id = response.json().id
        );
    }

    checkPercentage(): boolean {
        return this.selectedObject.weight + this.weight != 100;
    }
}

import {Component, Input, Output, EventEmitter, OnInit} from "@angular/core";
import {MainGroupService} from "main-group/service/main-group.service";
import {MainGroup} from "main-group/main-group";

@Component({
    selector: 'main-group-edit',
    templateUrl: '../ctemplate/html/ctemplate-edit.component.html',
    styleUrls: ['../ctemplate/ctemplate-edit.component.css'],
    providers: [MainGroupService]
})

export class MainGroupEditComponent implements OnInit {
    backupMainGroup: MainGroup;

    @Input() selectedObject: MainGroup;
    @Input() weight: number;
    @Output() onAction = new EventEmitter();

    constructor(private mainGroupService: MainGroupService) {

    }

    ngOnInit(): void {
        this.backupMainGroup = new MainGroup();
        this.backupMainGroup.name = this.selectedObject.name;
        this.backupMainGroup.weight = this.selectedObject.weight;
    }

    pressed(updatedMainGroup: MainGroup): void {
        if (updatedMainGroup) {
            this.updateMainGroup();
        } else {
            this.selectedObject.name = this.backupMainGroup.name;
            this.selectedObject.weight = this.backupMainGroup.weight;
        }
        this.onAction.emit();
    }

    updateMainGroup(): void {
        this.mainGroupService.updateMainGroup(this.selectedObject)
                .subscribe();
    }

    checkPercentage(): boolean {
        return this.selectedObject.weight + this.weight != 100;
    }

}

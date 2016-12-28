import {Component, Input, Output, EventEmitter, OnInit} from "@angular/core";
import {MainGroupService} from "main-group/service/main-group.service";
import {MainGroup} from "main-group/main-group";
import {AuditTemplate} from "../audit-template/audit-template";

@Component({
    selector: 'main-group-add',
    templateUrl: '../audit-template/html/audit-template-edit.component.html',
    styleUrls: ['../audit-template/audit-template-edit.component.css'],
    providers: [MainGroupService]
})

export class MainGroupAddComponent implements OnInit {
    selectedObject: MainGroup;

    @Input() auditTemplate: AuditTemplate;
    @Input() weight: number;
    @Output() onAdd = new EventEmitter<MainGroup>();

    constructor(private mainGroupService: MainGroupService) {

    }

    ngOnInit(): void {
        this.selectedObject = new MainGroup();
    }

    pressed(newMainGroup: MainGroup): void {
        if (newMainGroup) {
            this.addMainGroup();
        }
        this.onAdd.emit(newMainGroup);
    }

    addMainGroup(): void {
        this.selectedObject.idTemplate = this.auditTemplate.id;
        this.mainGroupService.setMainGroup(this.selectedObject).subscribe(
                response => this.selectedObject.id = response.json().id
        );
    }

    checkPercentage(): boolean {
        return this.selectedObject.weight + this.weight != 100;
    }
}

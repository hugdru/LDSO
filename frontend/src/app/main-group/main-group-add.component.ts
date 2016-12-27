import {Component, Input, Output, EventEmitter, OnInit} from "@angular/core";
import {MainGroupService} from "main-group/service/main-group.service";
import {MainGroup} from "main-group/main-group";
import {Ctemplate} from "../ctemplate/ctemplate";

@Component({
    selector: 'main-group-add',
    templateUrl: '../ctemplate/html/ctemplate-edit.component.html',
    styleUrls: ['../ctemplate/ctemplate-edit.component.css'],
    providers: [MainGroupService]
})

export class MainGroupAddComponent implements OnInit {
    selectedObject: MainGroup;

    @Input() ctemplate: Ctemplate;
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
        this.selectedObject.idTemplate = this.ctemplate.id;
        this.mainGroupService.setMainGroup(this.selectedObject).subscribe(
                response => this.selectedObject.id = response.json().id
        );
    }

    checkPercentage(): boolean {
        return this.selectedObject.weight + this.weight != 100;
    }
}

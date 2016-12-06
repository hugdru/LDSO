import {Component, Input, Output, EventEmitter, OnInit} from "@angular/core";
import {CtemplateService} from "ctemplate/service/ctemplate.service";
import {Ctemplate} from "ctemplate/ctemplate";

@Component({
    selector: 'ctemplate-edit',
    templateUrl: 'html/ctemplate-edit.component.html',
    styleUrls: ['ctemplate-edit.component.css'],
    providers: [CtemplateService]
})

export class CtemplateEditComponent implements OnInit {
    backupCtemplate: Ctemplate;

    @Input() objType: string;
    @Input() selectedObject: Ctemplate;
    @Output() onAction = new EventEmitter();

    constructor(private ctemplateService: CtemplateService) {

    }

    ngOnInit(): void {
        this.backupCtemplate = new Ctemplate();
        this.backupCtemplate.name = this.selectedObject.name;
    }

    pressed(updatedCtemplate: Ctemplate): void {
        if (updatedCtemplate) {
            this.updateCtemplate();
        } else {
            this.selectedObject.name = this.backupCtemplate.name;
        }
        this.onAction.emit();
    }

    updateCtemplate(): void {
        this.ctemplateService.updateCtemplate(this.selectedObject)
                .subscribe();
    }

}

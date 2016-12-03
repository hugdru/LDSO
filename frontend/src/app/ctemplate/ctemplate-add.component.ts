import {Component, Output, EventEmitter, OnInit} from "@angular/core";
import {CtemplateService} from "ctemplate/service/ctemplate.service";
import {Ctemplate} from "ctemplate/ctemplate";

@Component({
    selector: 'ctemplate-add',
    templateUrl: 'html/ctemplate-edit.component.html',
    styleUrls: ['ctemplate-edit.component.css'],
    providers: [CtemplateService]
})

export class CtemplateAddComponent implements OnInit {
    selectedObject: Ctemplate;

    @Output() onAdd = new EventEmitter<Ctemplate>();

    constructor(private ctemplateService: CtemplateService) {

    }

    ngOnInit(): void {
        this.selectedObject = new Ctemplate();
    }

    pressed(newCtemplate: Ctemplate): void {
        if (newCtemplate) {
            this.addCtemplate();
        }
        this.onAdd.emit(newCtemplate);
    }

    addCtemplate(): void {
        this.ctemplateService.setCtemplate(this.selectedObject).subscribe(
                response => this.selectedObject.id = response.json()
        );
    }
}

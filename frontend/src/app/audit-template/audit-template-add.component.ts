import {Component, Output, EventEmitter, OnInit, Input} from "@angular/core";
import {AuditTemplateService} from "./service/audit-template.service";
import {AuditTemplate} from "./audit-template";

@Component({
    selector: 'audit-template-add',
    templateUrl: 'html/audit-template-edit.component.html',
    styleUrls: ['audit-template-edit.component.css'],
    providers: [AuditTemplateService]
})

export class AuditTemplateAddComponent implements OnInit {
    selectedObject: AuditTemplate;

    @Input() objType: string;
    @Output() onAdd = new EventEmitter<AuditTemplate>();

    constructor(private auditTemplateService: AuditTemplateService) {

    }

    ngOnInit(): void {
        this.selectedObject = new AuditTemplate();
    }

    pressed(newAuditTemplate: AuditTemplate): void {
        if (newAuditTemplate) {
            this.addAuditTemplate();
        }
        this.onAdd.emit(newAuditTemplate);
    }

    addAuditTemplate(): void {
        this.auditTemplateService.setAuditTemplate(this.selectedObject).subscribe(
                response => this.selectedObject.id = response.json().id
        );
    }

}

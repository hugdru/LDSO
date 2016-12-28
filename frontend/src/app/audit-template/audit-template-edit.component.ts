import {Component, Input, Output, EventEmitter, OnInit} from "@angular/core";
import {AuditTemplateService} from "./service/audit-template.service";
import {AuditTemplate} from "./audit-template";

@Component({
    selector: 'auditTemplate-edit',
    templateUrl: 'html/audit-template-edit.component.html',
    styleUrls: ['audit-template-edit.component.css'],
    providers: [AuditTemplateService]
})

export class AuditTemplateEditComponent implements OnInit {
    backupAuditTemplate: AuditTemplate;

    @Input() objType: string;
    @Input() selectedObject: AuditTemplate;
    @Output() onAction = new EventEmitter();

    constructor(private auditTemplateService: AuditTemplateService) {

    }

    ngOnInit(): void {
        this.backupAuditTemplate = new AuditTemplate();
        this.backupAuditTemplate.name = this.selectedObject.name;
    }

    pressed(updatedAuditTemplate: AuditTemplate): void {
        if (updatedAuditTemplate) {
            this.updateAuditTemplate();
        } else {
            this.selectedObject.name = this.backupAuditTemplate.name;
        }
        this.onAction.emit();
    }

    updateAuditTemplate(): void {
        this.auditTemplateService.updateAuditTemplate(this.selectedObject)
                .subscribe();
    }

}

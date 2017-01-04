import {Component, OnInit} from "@angular/core";
import {AuditTemplateService} from "./service/audit-template.service";
import {AuditTemplate} from "./audit-template";
import {Close} from "./close";

@Component({
    selector: 'auditTemplate',
    templateUrl: 'html/audit-template.component.html',
    styleUrls: ['audit-template.component.css'],
    providers: [AuditTemplateService]
})

export class AuditTemplateComponent implements OnInit {
    auditTemplates: AuditTemplate[];
    parentAuditTemplate: AuditTemplate;
    objType: string;
    errorMsg: string;
    close: Close;



    constructor(private auditTemplateService: AuditTemplateService) {
        this.objType = "AuditTemplate"
    }

    ngOnInit(): void {
        this.initAuditTemplates();
    }

    initAuditTemplates(): void {
        this.auditTemplateService.getAuditTemplates().subscribe(
                data => this.auditTemplates = data,
                error => this.errorMsg = <any>error
        );
    }

    onDelete(auditTemplate: AuditTemplate): void {
        this.auditTemplateService.removeAuditTemplate(auditTemplate.id).subscribe();
        this.parentAuditTemplate = undefined;
    }

    onShow(auditTemplate: AuditTemplate): void {
        this.parentAuditTemplate = auditTemplate;
    }

    onClose(auditTemplate: AuditTemplate): void {
        this.close = new Close();
        this.close.close = true;
        this.auditTemplateService.closeAuditTemplate(
                auditTemplate.id, this.close).subscribe();
        this.parentAuditTemplate = undefined;
        this.initAuditTemplates();
    }

    onOpen(auditTemplate: AuditTemplate): void {
        this.close = new Close();
        this.close.close = false;
        this.auditTemplateService.closeAuditTemplate(
                auditTemplate.id, this.close).subscribe();
        this.parentAuditTemplate = undefined;
        this.initAuditTemplates();
    }

}

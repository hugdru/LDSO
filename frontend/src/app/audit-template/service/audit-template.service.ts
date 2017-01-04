import {Response} from "@angular/http";
import {Injectable} from "@angular/core";
import {Observable} from "rxjs/Observable";
import {
    auditTemplatesUrl, currentAuditTemplateUrl,
    closeAuditTemplateUrl
} from "shared/shared-data";
import {HandlerService} from "../../shared/service/handler.service";
import {AuditTemplate} from "../audit-template";
import {Close} from "../close";

@Injectable()
export class AuditTemplateService {

    constructor(private handler: HandlerService) {
    }

    getAuditTemplates(): Observable<AuditTemplate[]> {
        return this.handler.getAll<AuditTemplate[]>(auditTemplatesUrl);
    }

    getSomeAuditTemplates(tag: string, value: any): Observable<AuditTemplate[]> {
        return this.handler.getSome<AuditTemplate[]>(auditTemplatesUrl, tag, value);
    }

    getAuditTemplate(id: number): Observable<AuditTemplate> {
        return this.handler.get<AuditTemplate>(auditTemplatesUrl, id);
    }

    getCurrentAuditTemplate(): Observable<AuditTemplate> {
        return this.handler.getAll<AuditTemplate>(currentAuditTemplateUrl);
    }

    updateAuditTemplate(group: AuditTemplate): Observable<Response> {
        return this.handler.update<AuditTemplate>(auditTemplatesUrl, group, group.id);
    }

    setAuditTemplate(group: AuditTemplate): Observable<Response> {
        return this.handler.set<AuditTemplate>(auditTemplatesUrl, group);
    }

    removeAuditTemplate(id: number): Observable<Response> {
        return this.handler.delete(auditTemplatesUrl, id);
    }

    closeAuditTemplate(id: number, close: Close): Observable<Response> {
        return this.handler.set<Close>(
                closeAuditTemplateUrl.replace(/#/g, id.toString()), close);
    }

}

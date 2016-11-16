import { Response } from '@angular/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';

import { auditsUrl, auditsFindUrl } from 'shared/shared-data';
import { HandlerService } from 'handler.service';
import { Audit } from 'audit/audit';

@Injectable()
export class AuditService {

    constructor(private handler: HandlerService) { }

    getAudits(): Observable<Audit[]> {
        return this.handler.getAll<Audit[]>(auditsUrl);
    }

    getSomeAudits(tag: string, type: string, value: any)
    : Observable<Audit[]> {
        return this.handler.get<Audit[]>(auditsUrl, tag, type,
            value);
    }

    getAudit(tag: string, type: string, value: any): Observable<Audit> {
        return this.handler.get<Audit>(auditsFindUrl, tag, type, value);
    }

    updateAudit(audit: Audit): Observable<Response> {
        return this.handler.update<Audit>(auditsUrl, audit, audit._id);
    }

    setAudit(audit: Audit): Observable<Response> {
        return this.handler.set<Audit>(auditsUrl, audit);
    }

    removeAudit(id: number): Observable<Response> {
        return this.handler.delete(auditsUrl, id);
    }
}

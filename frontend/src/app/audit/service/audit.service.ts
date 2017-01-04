import {Response} from "@angular/http";
import {Injectable} from "@angular/core";
import {Observable} from "rxjs/Observable";
import {auditsUrl} from "shared/shared-data";
import {auditsSubGroupsUrl} from "shared/shared-data";
import {auditsCriterionUrl} from "shared/shared-data";
import {HandlerService} from "../../shared/service/handler.service";
import {Audit} from "audit/audit";
import {AuditCriterion} from "audit/audit";
import {SubGroup} from "sub-group/sub-group";

@Injectable()
export class AuditService {

    constructor(private handler: HandlerService) {
    }

    getAudits(): Observable<Audit[]> {
        return this.handler.getAll<Audit[]>(auditsUrl);
    }

    getSomeAudits(tag: string, value: any): Observable<Audit[]> {
        return this.handler.getSome<Audit[]>(auditsUrl, tag, value);
    }

    getAudit(id: number): Observable<Audit> {
        return this.handler.get<Audit>(auditsUrl, id);
    }

    updateAudit(audit: Audit): Observable<Response> {
        return this.handler.update<Audit>(auditsUrl, audit, audit.id);
    }

    setAudit(audit: Audit): Observable<Response> {
        return this.handler.set<Audit>(auditsUrl, audit);
    }

    removeAudit(id: number): Observable<Response> {
        return this.handler.delete(auditsUrl, id);
    }

	setAuditSubGroups(subGroups: SubGroup[]): Observable<Response> {
		return this.handler.set<SubGroup[]>(auditsSubGroupsUrl, subGroups);
	}

	setAuditCriterion(auditCriterion: AuditCriterion): Observable<Response> {
		return this.handler.set<AuditCriterion>(auditsCriterionUrl,
				auditCriterion);
	}
}

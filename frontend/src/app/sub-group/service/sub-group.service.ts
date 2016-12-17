import {Response} from "@angular/http";
import {Injectable} from "@angular/core";
import {Observable} from "rxjs/Observable";
import {loginUrl} from "shared/shared-data";
import {HandlerService} from "handler.service";
import {SubGroup} from "sub-group/sub-group";

@Injectable()
export class SubGroupService {
    constructor(private handler: HandlerService) {
    }

    getSubGroups(): Observable<SubGroup[]> {
        return this.handler.getAll<SubGroup[]>(loginUrl);
    }

    getSomeSubGroups(tag: string, value: any): Observable<SubGroup[]> {
        return this.handler.getSome<SubGroup[]>(loginUrl, tag, value);
    }

    getSubGroup(id: number): Observable<SubGroup> {
        return this.handler.get<SubGroup>(loginUrl, id);
    }

    updateSubGroup(subGroup: SubGroup): Observable<Response> {
        return this.handler.update<SubGroup>(loginUrl, subGroup,
                subGroup.id);
    }

    setSubGroup(subGroup: SubGroup): Observable<Response> {
        return this.handler.set<SubGroup>(loginUrl, subGroup);
    }

    removeSubGroup(id: number): Observable<Response> {
        return this.handler.delete(loginUrl, id);
    }
}

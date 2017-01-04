import {Response} from "@angular/http";
import {Injectable} from "@angular/core";
import {Observable} from "rxjs/Observable";
import {subGroupsUrl} from "shared/shared-data";
import {HandlerService} from "../../shared/service/handler.service";
import {SubGroup} from "sub-group/sub-group";

@Injectable()
export class SubGroupService {
    constructor(private handler: HandlerService) {
    }

    getSubGroups(): Observable<SubGroup[]> {
        return this.handler.getAll<SubGroup[]>(subGroupsUrl);
    }

    getSomeSubGroups(tag: string, value: any): Observable<SubGroup[]> {
        return this.handler.getSome<SubGroup[]>(subGroupsUrl, tag, value);
    }

    getSubGroup(id: number): Observable<SubGroup> {
        return this.handler.get<SubGroup>(subGroupsUrl, id);
    }

    updateSubGroup(subGroup: SubGroup): Observable<Response> {
        return this.handler.update<SubGroup>(subGroupsUrl, subGroup,
                subGroup.id);
    }

    setSubGroup(subGroup: SubGroup): Observable<Response> {
        return this.handler.set<SubGroup>(subGroupsUrl, subGroup);
    }

    removeSubGroup(id: number): Observable<Response> {
        return this.handler.delete(subGroupsUrl, id);
    }
}

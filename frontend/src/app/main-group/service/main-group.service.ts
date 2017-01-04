import {Response} from "@angular/http";
import {Injectable} from "@angular/core";
import {Observable} from "rxjs/Observable";
import {mainGroupsUrl} from "shared/shared-data";
import {HandlerService} from "../../shared/service/handler.service";
import {MainGroup} from "main-group/main-group";

@Injectable()
export class MainGroupService {

    constructor(private handler: HandlerService) {
    }

    getMainGroups(): Observable<MainGroup[]> {
        return this.handler.getAll<MainGroup[]>(mainGroupsUrl);
    }

    getSomeMainGroups(tag: string, value: any): Observable<MainGroup[]> {
        return this.handler.getSome<MainGroup[]>(mainGroupsUrl, tag, value);
    }

    getMainGroup(id: number): Observable<MainGroup> {
        return this.handler.get<MainGroup>(mainGroupsUrl, id);
    }

    updateMainGroup(group: MainGroup): Observable<Response> {
        return this.handler.update<MainGroup>(mainGroupsUrl, group, group.id);
    }

    setMainGroup(group: MainGroup): Observable<Response> {
        return this.handler.set<MainGroup>(mainGroupsUrl, group);
    }

    removeMainGroup(id: number): Observable<Response> {
        return this.handler.delete(mainGroupsUrl, id);
    }
}

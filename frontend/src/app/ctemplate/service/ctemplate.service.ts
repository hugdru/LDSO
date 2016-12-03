import {Response} from "@angular/http";
import {Injectable} from "@angular/core";
import {Observable} from "rxjs/Observable";
import {ctemplatesUrl} from "shared/shared-data";
import {HandlerService} from "handler.service";
import {Ctemplate} from "ctemplate/ctemplate";

@Injectable()
export class CtemplateService {

    constructor(private handler: HandlerService) {
    }

    getCtemplates(): Observable<Ctemplate[]> {
        return this.handler.getAll<Ctemplate[]>(ctemplatesUrl);
    }

    getSomeCtemplates(tag: string, value: any): Observable<Ctemplate[]> {
        return this.handler.getSome<Ctemplate[]>(ctemplatesUrl, tag, value);
    }

    getCtemplate(id: number): Observable<Ctemplate> {
        return this.handler.get<Ctemplate>(ctemplatesUrl, id);
    }

    updateCtemplate(group: Ctemplate): Observable<Response> {
        return this.handler.update<Ctemplate>(ctemplatesUrl, group, group.id);
    }

    setCtemplate(group: Ctemplate): Observable<Response> {
        return this.handler.set<Ctemplate>(ctemplatesUrl, group);
    }

    removeCtemplate(id: number): Observable<Response> {
        return this.handler.delete(ctemplatesUrl, id);
    }
}

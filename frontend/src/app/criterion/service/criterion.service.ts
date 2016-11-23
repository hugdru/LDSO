import {Response} from "@angular/http";
import {Injectable} from "@angular/core";
import {Observable} from "rxjs/Observable";
import {criteriaUrl, criteriaFindUrl} from "shared/shared-data";
import {HandlerService} from "handler.service";
import {Criterion} from "criterion/criterion";

@Injectable()
export class CriterionService {
    constructor(private handler: HandlerService) {
    }

    getCriteria(): Observable<Criterion[]> {
        return this.handler.getAll<Criterion[]>(criteriaUrl);
    }

    getSomeCriteria(tag: string, type: string, value: any): Observable<Criterion[]> {
        return this.handler.get<Criterion[]>(criteriaUrl, tag, type,
                value);
    }

    getCriterion(tag: string, type: string, value: any): Observable<Criterion> {
        return this.handler.get<Criterion>(criteriaFindUrl, tag, type,
                value);
    }

    updateCriterion(criterion: Criterion): Observable<Response> {
        return this.handler.update<Criterion>(criteriaUrl, criterion,
                criterion._id);
    }

    setCriterion(criterion: Criterion): Observable<Response> {
        return this.handler.set<Criterion>(criteriaUrl, criterion);
    }

    removeCriterion(id: number): Observable<Response> {
        return this.handler.delete(criteriaUrl, id);
    }

}

import {Response} from "@angular/http";
import {Injectable} from "@angular/core";
import {Observable} from "rxjs/Observable";
import {criteriaUrl} from "shared/shared-data";
import {HandlerService} from "../../shared/service/handler.service";
import {Criterion} from "criterion/criterion";

@Injectable()
export class CriterionService {
    constructor(private handler: HandlerService) {
    }

    getCriteria(): Observable<Criterion[]> {
        return this.handler.getAll<Criterion[]>(criteriaUrl);
    }

    getSomeCriteria(tag: string, value: any): Observable<Criterion[]> {
        return this.handler.getSome<Criterion[]>(criteriaUrl, tag, value);
    }

    getCriterion(id: number): Observable<Criterion> {
        return this.handler.get<Criterion>(criteriaUrl, id);
    }

    updateCriterion(criterion: Criterion): Observable<Response> {
        return this.handler.update<Criterion>(criteriaUrl, criterion,
                criterion.id);
    }

    setCriterion(criterion: Criterion): Observable<Response> {
        return this.handler.set<Criterion>(criteriaUrl, criterion);
    }

    removeCriterion(id: number): Observable<Response> {
        return this.handler.delete(criteriaUrl, id);
    }

}

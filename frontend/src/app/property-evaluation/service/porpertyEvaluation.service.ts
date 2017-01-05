import {Response} from "@angular/http";
import {Injectable} from "@angular/core";
import {HandlerService} from "../../shared/service/handler.service";
import {Observable} from "rxjs/Observable";
import {PropertyEvaluation} from "../property-evaluation";



@Injectable()
export class PropertyEvaluationService {

    constructor(private handler: HandlerService) {
    }
    getPropertyEvaluation(id: number): Observable<PropertyEvaluation[]> {
        let formated = "/audits?idProperty=" + id ;
        return this.handler.getAll<PropertyEvaluation[]>(formated);
    }


}

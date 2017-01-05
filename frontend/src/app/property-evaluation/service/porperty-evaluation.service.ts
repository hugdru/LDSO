import {Injectable} from "@angular/core";
import {HandlerService} from "../../handler.service";
import {Observable} from "rxjs/Observable";
import {PropertyEvaluation} from "../property-evaluation";
import {auditsUrl} from "../../shared/shared-data";



@Injectable()
export class PropertyEvaluationService {

    constructor(private handler: HandlerService) {
    }

    getSomePropertyEvaluation(tag: string, value: any): Observable<PropertyEvaluation[]> {
        return this.handler.getSome<PropertyEvaluation[]>(auditsUrl, tag, value);
    }

}

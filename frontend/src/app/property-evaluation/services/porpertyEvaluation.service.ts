import {Response} from "@angular/http";
import {Injectable} from "@angular/core";
import {Observable} from "rxjs/Observable";
import {propertiesUrl} from "shared/shared-data";
import {HandlerService} from "handler.service";
import {PropertyEvaluation} from "property-evaluation/property-evaluation";

@Injectable()
export class PropertyEvaluationService {

    constructor(private handler: HandlerService) {

    }
    getPropertyEvaluation(id: number): Observable<PropertyEvaluation> {
        return this.handler.get<PropertyEvaluation>(propertiesUrl, id);
    }
}

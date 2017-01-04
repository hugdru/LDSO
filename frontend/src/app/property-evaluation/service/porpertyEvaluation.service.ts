import {Response} from "@angular/http";
import {Injectable} from "@angular/core";
import {HandlerService} from "../../handler.service";
import {Observable} from "rxjs/Observable";


@Injectable()
export class PropertyEvaluationService {

    constructor(private handler: HandlerService) {
    }

    getPropertyEvaluation(id: number): Observable<Response> {
        let formated = "/properties/" + id + "/audits";
        return this.handler.getResponse(formated);
    }


}

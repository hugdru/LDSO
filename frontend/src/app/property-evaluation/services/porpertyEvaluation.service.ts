import {Http, Response, Headers, RequestOptions} from "@angular/http";
import {Injectable} from "@angular/core";
import {PropertyEvaluation} from "property-evaluation/property-evaluation";
import 'rxjs/add/observable/throw';
import "rxjs/add/operator/map";
import "rxjs/add/operator/catch";

@Injectable()
export class PropertyEvaluationService {

    headers = new Headers({ 'Content-Type': 'application/json' });
    options = new RequestOptions({ headers: this.headers });

    constructor(private http: Http) {
    }

    getPropertyEvaluation(id: number): PropertyEvaluation {
        let formated = "/properties/" + id + "/audits";
       return null;

    }



}

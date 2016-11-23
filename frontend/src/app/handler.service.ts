import {Injectable} from "@angular/core";
import {Http, Response} from "@angular/http";
import {Observable} from "rxjs/Observable";
import "rxjs/add/operator/map";
import "rxjs/add/operator/catch";

@Injectable()
export class HandlerService {

    constructor(private http: Http) {
    }

    private handleError(error: Response | any) {
        let errMsg: string;

        if (error instanceof Response) {
            const body = error.json() || '';
            const err = body.error || JSON.stringify(body);
            errMsg = `${error.status} - ${error.statusText || ''} ${err}`;
        } else {
            errMsg = error.message ? error.message : error.toString();
        }

        console.error(errMsg);
        return Observable.throw(errMsg);
    }

    get<T>(url: string, tag: string, type: string, value: any): Observable<T> {
        let formated = url + "?tag=" + tag + "&type=" + type
                + "&value=" + value;
        return this.getAll<T>(formated);
    }

    getAll<T>(url: string): Observable<T> {
        return this.http.get(url)
                .map((response: Response) => response.json())
                .map((data: any) => {
                    let result: T = null;
                    if (data) {
                        result = data;
                    }
                    return result;
                }).catch(this.handleError);
    }

    update<T>(url: string, object: T, id: number): Observable<Response> {
        return this.http.put(url + "?_id=" + id, JSON.stringify(object))
                .map((result: Response) => result);
    }

    delete(url: string, id: number): Observable<Response> {
        let formated = url + "?_id=" + id;
        return this.http.delete(formated).map((result: Response) => result);
    }

    set<T>(url: string, object: T): Observable<Response> {
        return this.http.post(url, JSON.stringify(object))
                .map((response: Response) => response);
    }

}

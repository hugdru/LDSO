import {Injectable} from "@angular/core";
import {Http, Response, Headers, RequestOptions} from "@angular/http";
import {Observable} from "rxjs/Observable";
import 'rxjs/add/observable/throw';
import "rxjs/add/operator/map";
import "rxjs/add/operator/catch";

@Injectable()
export class HandlerService {

    headers = new Headers({ 'Content-Type': 'application/json' });
    options = new RequestOptions({ headers: this.headers });

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

    get<T>(url: string, id: number): Observable<T> {
        let formated = url + "/" + id;
        return this.getAll<T>(formated);
    }

    getSome<T>(url: string, tag: string, value: any): Observable<T> {
        let formated = url + "?" + tag + "=" + value;
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
        return this.http.put(url + "/" + id, JSON.stringify(object), this.options)
                .map((response: Response) => response)
                .catch(this.handleError);
    }

    delete(url: string, id: number): Observable<Response> {
        let formated = url + "/" + id;
        return this.http.delete(formated)
                .map((response: Response) => response)
                .catch(this.handleError);
    }

    set<T>(url: string, object: T): Observable<Response> {
        return this.http.post(url, JSON.stringify(object), this.options)
                .map((response: Response) => response)
                .catch(this.handleError);
    }

    login<T>(url: string, object: T): Observable<Response> {
        return this.http.post(url, JSON.stringify(object), this.options)
                .map((response: Response) => response.json())
                .map((response) => {
                    if (response.success) {
                        localStorage.setItem('auth_token', response.auth_token);
                    }
                    return response.success;
                })
                .catch(this.handleError);
    }



}

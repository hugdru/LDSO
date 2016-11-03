import { Injectable } from '@angular/core';
import { Http, Response } from '@angular/http';
import { Observable } from 'rxjs/Observable';

import 'rxjs/add/operator/map';

@Injectable()
export class HandlerService {

	constructor(private http: Http) {}

	private handleError (error: Response | any) {
		let errMsg: string;

		if(error instanceof Response) {
			const body = error.json() || '';
			const err = body.error || JSON.stringify(body);
				errMsg = `${error.status} - ${error.statusText || ''} ${err}`;
			}
			else {
				errMsg = error.message ? error.message : error.toString();
			}

			console.error(errMsg);

			return Observable.throw(errMsg);
		}

		get<T>(url: string): Observable<T> {
			return this.http.get(url)
					.map((result: Response) => result.json())
					.map((data: any) => {
						let result: T = null;
						if(data) {
							result = data;
						}
						return result;
					}).catch(this.handleError);
		}

		set<T>(object: T, url: string): void {
			this.http.post(url, JSON.stringify(object))
					.map(res => res.json());
		}
}

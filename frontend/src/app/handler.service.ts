import { Injectable } from '@angular/core';
import { Http, Response } from '@angular/http';
import { Observable } from 'rxjs/Observable';

import 'rxjs/add/operator/map';

@Injectable()
export class HandlerService {

	constructor(private http: Http) {}

	get<T>(url: string): Observable<T> {
		return this.http.get(url)
				.map((result: Response) => result.json())
				.map((data: any) => {
					let result: T = null;
					if(data) {
						result = data;
					}
					return result;
				});
	}
}

import { Injectable } from '@angular/core';
import { Http, Response } from '@angular/http';
import { Observable } from 'rxjs/Observable';

import 'rxjs/add/operator/map';

import { Criterion } from './criterion';

@Injectable()
export class CriteriaService {
	url = 'http://localhost:8080';
	urlGetCrit = this.url + '/criteria';
	criteria: Criterion[];

	constructor(private http: Http) {}

	getCriteria(): Observable<Criterion[]> {
		return this.http.get(this.urlGetCrit )
				.map((result: Response) => result.json())
				.map((data: any) => {
					 let result: Criterion[] = null;
					 if(data) {
						this.criteria = data;
					 }
					return this.criteria;});
	}
}

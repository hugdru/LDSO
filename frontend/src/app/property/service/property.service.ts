import { Injectable } from '@angular/core';
import { Http, Response } from '@angular/http';
import { Observable } from 'rxjs/Observable';

import 'rxjs/add/operator/map';

import { Property } from './property';

@Injectable()
export class PropertyService {
	prop = 'Hotel Sunny';
	// property: Observable<Property>;
	property: Property;

	url = 'http://localhost:8080';
	urlGetProp = this.url + '/property';

	constructor(private http: Http) {}


	getProperty(prop: string): Observable<Property> {
		return this.http.get(this.urlGetProp + '?label=name&value=' + prop)
				.map((result: Response) => result.json())
				.map((data: any) => {
					 let result: Property = null;
					 if(data) {
						this.property = data;
					 }
					return this.property;});
	}
}

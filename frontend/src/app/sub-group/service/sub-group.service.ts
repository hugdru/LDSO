import { Injectable } from '@angular/core';
import { Http, Response } from '@angular/http';
import { Observable } from 'rxjs/Observable';

import 'rxjs/add/operator/map';

import { SharedData } from 'shared/shared-data';

@Injectable()
export class SubGroupService {
	private subGroups: SubGroup[];

	 constructor(private http: Http) { }

	 getSubGroup(name: string): Observable<SubGroup> {
	 }
}

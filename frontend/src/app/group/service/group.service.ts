import { Injectable } from '@angular/core';
import { Http, Response } from '@angular/http';
import { Observable } from 'rxjs/Observable';

import 'rxjs/add/operator/map';

import { getGroupUrl, setGroupUrl } from 'shared/shared-data';
import { Group } from 'group/group';

@Injectable()
export class GroupService {
	private groups: Group[];

	constructor(private http: Http) {}

	getGroups(): Observable<Group[]> {
		return this.http.get(getGroupUrl)
				.map((result: Response) => result.json())
				.map((data: any) => {
					let result: Group[] = null;
					if(data) {
						result = data;
					}
					return result;
				});
	}

	setGroup(group: Group): void {
		this.http.post(setGroupUrl, JSON.stringify(group))
				.map(res => res.json()).subscribe();
	}
}

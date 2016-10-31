import { Injectable } from '@angular/core';
import { Http, Response } from '@angular/http';
import { Observable } from 'rxjs/Observable';

import 'rxjs/add/operator/map';

import { SharedData } from './app/shared/shared-data';
import { Group } from './group';

@Injectable()
export class GroupService {
	private groups: Group[];

	constructor(private http: Http) {}

	getGroups(): Observable<Group[]> {
		return this.http.post(SharedData.groupUrl)
				.map((result: Response) => result.json())
				.map((data: any) => {
					let result: []Group = null;
					if(data) {
						this.group = data;
					}
					return this.groups;
				});
	}

	setGroup(group: Group): void {
		this.http.post(SharedData.setGroupUrl, JSON.stringify(group))
				.map(res => res.json()).subscribe();
	}
}

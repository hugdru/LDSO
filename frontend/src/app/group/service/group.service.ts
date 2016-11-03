import { Injectable } from '@angular/core';
import { Http, Response } from '@angular/http';
import { Observable } from 'rxjs/Observable';

import 'rxjs/add/operator/map';

import { getGroupUrl, setGroupUrl } from 'shared/shared-data';
import { HandlerService } from 'handler.service';
import { Group } from 'group/group';

@Injectable()
export class GroupService {
	constructor(private http: Http, private handler: HandlerService) {}

	getGroups(): Observable<Group[]> {
		return this.handler.get<Group[]>(getGroupUrl);
	}

	setGroup(group: Group): void {
		this.http.post(setGroupUrl, JSON.stringify(group))
				.map(res => res.json()).subscribe();
	}
}

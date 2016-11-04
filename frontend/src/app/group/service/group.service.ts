import { Response } from '@angular/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';

import {
	getGroupsUrl,
	getGroupUrl,
	updateGroupUrl,
	setGroupUrl,
	removeGroupUrl
} from 'shared/shared-data';
import { HandlerService } from 'handler.service';
import { Group } from 'group/group';

@Injectable()
export class GroupService {
	constructor(private handler: HandlerService) { }

	getGroups(): Observable<Group[]> {
		return this.handler.get<Group[]>(getGroupsUrl);
	}

	getGroup(tag: string, type: string, value: any): Observable<Group> {
		return this.handler.get<Group>(getGroupUrl + "?tag=" + tag + "&type="
			   + type + "&value=" + value);
	}

	updateGroup(id: number, tag: string, type: string, value: any)
			: Observable<Response> {
		return this.handler.change(updateGroupUrl + "?_id"+ id + "&tag=" + tag 
			   + "&type=" + type + "&value=" + value);
	}

	setGroup(group: Group): Observable<Response> {
		return this.handler.set<Group>(group, setGroupUrl);
	}

	removeGroup(id: number): Observable<Response> {
		return this.handler.change(removeGroupUrl + "?_id" + id);
	}
}

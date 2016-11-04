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
		return this.handler.getOne<Group>(getGroupUrl, tag, type, value);
	}

	updateGroup(id: number, tag: string, type: string, value: any)
			: Observable<Response> {
		return this.handler.update(updateGroupUrl, id, tag, type, value);
	}

	setGroup(group: Group): Observable<Response> {
		return this.handler.set<Group>(group, setGroupUrl);
	}

	removeGroup(id: number): Observable<Response> {
		return this.handler.delete(removeGroupUrl, id);
	}
}

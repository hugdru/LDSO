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

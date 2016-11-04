import { Response } from '@angular/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';

import {
	getMainGroupsUrl,
	getMainGroupUrl,
	updateMainGroupUrl,
	setMainGroupUrl,
	removeMainGroupUrl
} from 'shared/shared-data';
import { HandlerService } from 'handler.service';
import { MainGroup } from 'main-group/main-group';

@Injectable()
export class MainGroupService {

	constructor(private handler: HandlerService) { }

	getMainGroups(): Observable<MainGroup[]> {
		return this.handler.get<MainGroup[]>(getMainGroupsUrl);
	}

	getMainGroup(tag: string, type: string, value: any): Observable<MainGroup> {
		return this.handler.getOne<MainGroup>(getMainGroupUrl, tag, type, value);
	}

	updateMainGroup(id: number, tag: string, type: string, value: any)
			: Observable<Response> {
		return this.handler.update(updateMainGroupUrl, id, tag, type, value);
	}

	setMainGroup(group: MainGroup): Observable<Response> {
		return this.handler.set<MainGroup>(group, setMainGroupUrl);
	}

	removeMainGroup(id: number): Observable<Response> {
		return this.handler.delete(removeMainGroupUrl, id);
	}
}

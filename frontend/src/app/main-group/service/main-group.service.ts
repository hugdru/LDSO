import { Response } from '@angular/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';

import { mainGroupsUrl, mainGroupsFindUrl } from 'shared/shared-data';
import { HandlerService } from 'handler.service';
import { MainGroup } from 'main-group/main-group';

@Injectable()
export class MainGroupService {

	constructor(private handler: HandlerService) { }

	getMainGroups(): Observable<MainGroup[]> {
		return this.handler.getAll<MainGroup[]>(mainGroupsUrl);
	}

	getSomeMainGroups(tag: string, type: string, value: any)
			: Observable<MainGroup[]> {
		return this.handler.get<MainGroup[]>(mainGroupsUrl, tag, type,
				 value);
	}

	getMainGroup(tag: string, type: string, value: any): Observable<MainGroup> {
		return this.handler.get<MainGroup>(mainGroupsFindUrl, tag, type, value);
	}

	updateMainGroup(group: MainGroup): Observable<Response> {
		return this.handler.update<MainGroup>(mainGroupsUrl, group, group._id);
	}

	setMainGroup(group: MainGroup): Observable<Response> {
		return this.handler.set<MainGroup>(mainGroupsUrl, group);
	}

	removeMainGroup(id: number): Observable<Response> {
		return this.handler.delete(mainGroupsUrl, id);
	}
}

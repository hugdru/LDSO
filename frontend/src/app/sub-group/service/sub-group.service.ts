import { Response } from '@angular/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';

import { subGroupsUrl, subGroupsFindUrl } from 'shared/shared-data';
import { HandlerService } from 'handler.service';
import { SubGroup } from 'sub-group/sub-group';

@Injectable()
export class SubGroupService {
	constructor(private handler: HandlerService) { }

	getSubGroups(): Observable<SubGroup[]> {
	    return this.handler.getAll<SubGroup[]>(subGroupsUrl);
	}

	getSomeSubGroups(tag: string, type: string, value: any)
			: Observable<SubGroup[]> {
	    return this.handler.get<SubGroup[]>(subGroupsUrl, tag, type, value);
	}

	getSubGroup(tag: string, type: string, value: any): Observable<SubGroup> {
	    return this.handler.get<SubGroup>(subGroupsFindUrl, tag, type, value);
	}

	updateSubGroup(subGroup: SubGroup): Observable<Response> {
		return this.handler.update<SubGroup>(subGroupsUrl, subGroup);
	}

	setSubGroup(subGroup: SubGroup): Observable<Response> {
	    return this.handler.set<SubGroup>(subGroupsUrl, subGroup);
	}

	removeSubGroup(id: number): Observable<Response> {
		return this.handler.delete(subGroupsUrl, id);
	}
}

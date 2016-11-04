import { Response } from '@angular/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';

import {
	getSubGroupsUrl,
	getSubGroupUrl,
	setSubGroupUrl
} from 'shared/shared-data';
import { HandlerService } from 'handler.service';
import { SubGroup } from 'sub-group/sub-group';

@Injectable()
export class SubGroupService {
	 constructor(private handler: HandlerService) { }

	 getSubGroups(): Observable<SubGroup[]> {
		 return this.handler.get<SubGroup[]>(getSubGroupsUrl);
	 }

	 getSubGroup(tag: string, type: string, value: any): Observable<SubGroup> {
		 return this.handler.get<SubGroup>(getSubGroupUrl + "?tag=" + tag
			   + "&type=" + type + "&value=" + value);
	 }

	 setSubGroup(subGroup: SubGroup): Observable<Response> {
		 return this.handler.set<SubGroup>(subGroup, setSubGroupUrl);
	 }
}

import { Response } from '@angular/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';

import {
	getCriteriaUrl,
	getCriterionUrl,
	setCriterionUrl
} from 'shared/shared-data';
import { HandlerService } from 'handler.service';
import { Criterion } from 'criterion/criterion';

@Injectable()
export class CriterionService {
	constructor(private handler: HandlerService) {}

	getCriteria(): Observable<Criterion[]> {
		return this.handler.get<Criterion[]>(getCriteriaUrl);
	}

	getCriterion(tag: string, type: string, value: any)
			: Observable<Criterion> {
		return this.handler.get<Criterion>(getCriterionUrl + "?tag=" + tag
				+ "&type=" + type + "&value=" + value);
	}

	setCriterion(criterion: Criterion): Observable<Response> {
		return this.handler.set<Criterion>(criterion, setCriterionUrl);
	}

}

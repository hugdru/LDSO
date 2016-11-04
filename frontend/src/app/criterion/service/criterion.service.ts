import { Response } from '@angular/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';

import {
	getCriteriaUrl,
	getSomeCriteriaUrl,
	getCriterionUrl,
	updateCriterionUrl,
	setCriterionUrl,
	removeCriterionUrl
} from 'shared/shared-data';
import { HandlerService } from 'handler.service';
import { Criterion } from 'criterion/criterion';

@Injectable()
export class CriterionService {
	constructor(private handler: HandlerService) {}

	getCriteria(): Observable<Criterion[]> {
		return this.handler.getAll<Criterion[]>(getCriteriaUrl);
	}

	getSomeCriteria(tag: string, type: string, value: any)
			: Observable<Criterion[]> {
		return this.handler.get<Criterion[]>(getCriterionUrl, tag, type,
				value);
	}

	getCriterion(tag: string, type: string, value: any)
			: Observable<Criterion> {
		return this.handler.getOne<Criterion>(getCriterionUrl, tag, type,
				value);
	}

	updateCriterion(id: number, tag: string, type: string, value: any)
			:Observable<Response> {
		return this.handler.update(updateCriterionUrl, id, tag, type, value);
	}

	setCriterion(criterion: Criterion): Observable<Response> {
		return this.handler.set<Criterion>(criterion, setCriterionUrl);
	}

	removeCriterion(id: number): Observable<Response> {
		return this.handler.delete(removeCriterionUrl, id);
	}

}

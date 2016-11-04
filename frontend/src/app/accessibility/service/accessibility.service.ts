import { Response } from '@angular/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';

import {
	getAccessibilitiesUrl,
	getSomeAccessibilitiesUrl,
	getAccessibilityUrl,
	updateAccessibilityUrl,
	setAccessibilityUrl,
	removeAccessibilityUrl
} from 'shared/shared-data';
import { HandlerService } from 'handler.service';
import { Accessibility } from 'accessibility/accessibility';

@Injectable()
export class AccessibilityService {
	constructor(private handler: HandlerService) { }

	getAccessibilities(): Observable<Accessibility[]> {
		return this.handler.getAll<Accessibility[]>(getAccessibilitiesUrl);
	}

	getSomeAccessibilities(tag: string, type: string, value: any)
			: Observable<Accessibility[]> {
		return this.handler.get<Accessibility[]>(getSomeAccessibilitiesUrl,
				tag, type, value);
	}

	getAccessibility(tag: string, type: string, value: any)
			: Observable<Accessibility> {
		return this.handler.get<Accessibility>(getAccessibilitiesUrl, tag, type,
				value);
	}

	updateAccessibility(id: number, tag: string, type: string, value: any)
			:Observable<Response> {
		return this.handler.update(updateAccessibilityUrl, id, tag, type,
				value);
	}

	setAccessibility(accessibility: Accessibility): Observable<Response> {
		return this.handler.set<Accessibility>(accessibility,
				setAccessibilityUrl);
	}

	removeAccessibility(id: number): Observable<Response> {
		return this.handler.delete(removeAccessibilityUrl, id);
	}

}

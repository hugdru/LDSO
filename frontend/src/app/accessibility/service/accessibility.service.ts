import { Response } from '@angular/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';

import {
	getAccessibilitiesUrl,
	getAccessibilityUrl,
	setAccessibilityUrl
} from 'shared/shared-data';
import { HandlerService } from 'handler.service';
import { Accessibility } from 'accessibility/accessibility';

@Injectable()
export class AccessibilityService {

	constructor(private handler: HandlerService) { }

	getAccessibilities(): Observable<Accessibility[]> {
		return this.handler.get<Accessibility[]>(getAccessibilitiesUrl);
	}

	getAccessibility(tag: string, type: string, value: any)
			: Observable<Accessibility> {
		return this.handler.get<Accessibility>(getAccessibilitiesUrl + "?tag="
				+ tag + "&type=" + type + "&value=" + value);
	}

	setAccessibility(accessibility: Accessibility): Observable<Response> {
		return this.handler.set<Accessibility>(accessibility,
				setAccessibilityUrl);
	}

}

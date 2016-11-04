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

	getAccessibility(tag: string, type: string, value: string)
			: Observable<Accessibility> {
		return this.handler.get<Accessibility>(getAccessibilityUrl
			   + "?tag=" + tag + "&type=" + type + "&value=" + value);
	}

	// setAccessibility(accessibility: Accessibility): void {
	// 	this.handler.set<Accessibility>(group, setAccessibilityUrl);
	// }
}

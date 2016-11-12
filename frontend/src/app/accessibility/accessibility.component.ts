import {
	Component,
	OnInit,
	OnChanges,
	SimpleChanges,
	Input
} from '@angular/core';

import {
	AccessibilityService
} from 'accessibility/service/accessibility.service';
import { Accessibility } from 'accessibility/accessibility';
import { Criterion } from 'criterion/criterion';

@Component({
	selector: 'accessibility',
	templateUrl: 'html/accessibility.component.html',
	styleUrls: [ '../main-group/main-group.component.css' ],
	providers: [ AccessibilityService ]
})

export class AccessibilityComponent implements OnInit, OnChanges {
	accessibilities: Accessibility[];

	@Input() parentCriterion: Criterion;

	constructor(private accessibilityService: AccessibilityService){ }

	ngOnChanges(changes: SimpleChanges): void {
		for(let i in changes) {
			this.initAccessibilities(changes[i].currentValue._id);
		}
	}

	ngOnInit() {
		this.initAccessibilities(this.parentCriterion._id);
	}

	initAccessibilities(criterionId: number): void {
		this.accessibilityService
			.getSomeAccessibilities("criterion", "int", criterionId)
			.subscribe(data => this.accessibilities = data);

	}

	onDelete(accessibility: Accessibility): void {
		this.accessibilityService.removeAccessibility(accessibility._id)
				.subscribe();
	}

	onShow(accessibility: Accessibility): void { }

}

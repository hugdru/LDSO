import {
	Component,
	OnInit,
	OnChanges,
	SimpleChanges,
	Input
} from '@angular/core';

import { CriterionService } from 'criterion/service/criterion.service';
import { Criterion } from 'criterion/criterion';
import { SubGroup } from 'sub-group/sub-group';

@Component({
	selector: 'criterion',
	templateUrl: 'html/criterion.component.html',
	styleUrls: [ '../main-group/main-group.component.css' ],
	providers: [ CriterionService ]
})

export class CriterionComponent implements OnInit, OnChanges {
	criteria: Criterion[];
	selectedShowAccessibilities: Criterion;

	@Input() selectedShowCriterion: SubGroup;

	constructor(private criterionService: CriterionService){ }

	ngOnChanges(changes: SimpleChanges): void {
		for(let i in changes) {
			console.log(changes[i].previousValue);
			console.log(changes[i].currentValue);
			this.initCriteria(changes[i].currentValue._id);
			console.log(this.criteria);
		}
	}

	ngOnInit() {
		this.initCriteria(this.selectedShowCriterion._id);
		console.log(this.selectedShowCriterion);
		console.log(this.criteria);
	}

	initCriteria(subGroupId: number): void {
		this.criterionService
			.getSomeCriteria("sub_group", "int", subGroupId)
			.subscribe(data => this.criteria = data);

	}

	onDelete(criterion: Criterion): void {
		this.criterionService.removeCriterion(criterion._id).subscribe();
	}

	onShow(criterion: Criterion): void {
		this.selectedShowAccessibilities = criterion;
	}

}

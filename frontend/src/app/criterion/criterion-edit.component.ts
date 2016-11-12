import {
	Component,
	Input,
	Output,
	EventEmitter,
	OnInit
} from '@angular/core';

import { CriterionService } from 'criterion/service/criterion.service';
import { Criterion } from 'criterion/criterion';

@Component({
	selector: 'criterion-edit',
	templateUrl: '../main-group/html/main-group-edit.component.html',
	styleUrls: [ '../main-group/main-group-edit.component.css' ],
	providers: [ CriterionService ]
})

export class CriterionEditComponent implements OnInit {
	backupCriterion: Criterion;

	@Input() objType: string;
	@Input() selectedObject: Criterion;
	@Input() weight: number;
	@Output() onAction = new EventEmitter();

	constructor(private criterionService: CriterionService) {

	}

	ngOnInit(): void {
		console.log(this.objType);
		this.backupCriterion = new Criterion();
		this.backupCriterion.name = this.selectedObject.name;
		this.backupCriterion.weight = this.selectedObject.weight;
	}

	pressed(updatedCriterion: Criterion): void {
		if(updatedCriterion) {
			this.updateCriterion();
		} else {
			this.selectedObject.name = this.backupCriterion.name;
			this.selectedObject.weight = this.backupCriterion.weight;
		}
		this.onAction.emit();
	}

	updateCriterion(): void {
		this.criterionService.updateCriterion(this.selectedObject)
				.subscribe();
	}

	checkPercentage(): boolean {
		return this.selectedObject.weight + this.weight > 100;
	}

}

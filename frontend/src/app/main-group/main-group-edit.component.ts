import {
	Component,
	Input,
	Output,
	EventEmitter,
	OnInit
} from '@angular/core';

import { MainGroupService } from 'main-group/service/main-group.service';
import { MainGroup } from 'main-group/main-group';

@Component({
	selector: 'main-group-edit',
	templateUrl: 'html/main-group-edit.component.html',
	styleUrls: [ 'main-group-edit.component.css' ],
	providers: [ MainGroupService ]
})

export class MainGroupEditComponent implements OnInit {
	backupMainGroup: MainGroup;

	@Input() selectedMainGroup: MainGroup;
	@Input() weight: number;
	@Output() onAction = new EventEmitter();

	constructor(private mainGroupService: MainGroupService) {

	}

	ngOnInit(): void {
		this.backupMainGroup = new MainGroup();
		this.backupMainGroup.name = this.selectedMainGroup.name;
		this.backupMainGroup.weight = this.selectedMainGroup.weight;
	}

	pressed(updatedMainGroup: MainGroup): void {
		if(updatedMainGroup) {
			this.updateMainGroup();
		} else {
			this.selectedMainGroup.name = this.backupMainGroup.name;
			this.selectedMainGroup.weight = this.backupMainGroup.weight;
		}
		this.onAction.emit();
	}

	updateMainGroup(): void {
		this.mainGroupService.updateMainGroup(this.selectedMainGroup)
				.subscribe();
	}

	checkPercentage(): boolean {
		return this.selectedMainGroup.weight + this.weight > 100;
	}

}

import { Component, Input, Output, EventEmitter } from '@angular/core';

import { MainGroupService } from 'main-group/service/main-group.service';
import { MainGroup } from 'main-group/main-group';

@Component({
	selector: 'main-group-edit',
	templateUrl: 'html/main-group-edit.component.html',
	providers: [ MainGroupService ]
})

export class MainGroupEditComponent {
	@Input() selectedMainGroup: MainGroup;
	@Output() onAction = new EventEmitter<MainGroup>();

	constructor(private mainGroupService: MainGroupService) {

	}

	pressed(update: boolean, updatedMainGroup: MainGroup): void {
		if(update) {
			this.updateMainGroup();
		}
		this.onAction.emit(updatedMainGroup);
	}

	updateMainGroup(): void {
		this.mainGroupService.updateMainGroup(this.selectedMainGroup)
				.subscribe();
	}

}

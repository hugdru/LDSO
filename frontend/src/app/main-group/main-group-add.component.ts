import { Component, Output, EventEmitter, OnInit } from '@angular/core';

import { MainGroupService } from 'main-group/service/main-group.service';
import { MainGroup } from 'main-group/main-group';

@Component({
	selector: 'main-group-add',
	templateUrl: 'html/main-group-edit.component.html',
	providers: [ MainGroupService ]
})

export class MainGroupAddComponent implements OnInit {
	selectedMainGroup: MainGroup;

	@Output() onAdd = new EventEmitter<MainGroup>();

	constructor(private mainGroupService: MainGroupService) {

	}

	ngOnInit(): void {
		this.selectedMainGroup = new MainGroup();
	}

	pressed(newMainGroup: MainGroup): void {
		if(newMainGroup) {
			this.addMainGroup();
		}
		this.onAdd.emit(newMainGroup);
	}

	addMainGroup(): void {
		this.mainGroupService.setMainGroup(this.selectedMainGroup).subscribe();
	}
}

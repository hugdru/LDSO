import { Component, Input, Output, EventEmitter, OnInit } from '@angular/core';

import { MainGroupService } from 'main-group/service/main-group.service';
import { MainGroup } from 'main-group/main-group';

@Component({
	selector: 'main-group-add',
	templateUrl: 'html/main-group-edit.component.html',
	styleUrls: [ 'main-group-edit.component.css' ],
	providers: [ MainGroupService ]
})

export class MainGroupAddComponent implements OnInit {
	selectedMainGroup: MainGroup;

	@Input() weight: number;
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

	checkPercentage(): boolean {
		return this.selectedMainGroup.weight + this.weight > 100;
	}
}

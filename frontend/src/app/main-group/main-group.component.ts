import { Component, OnInit, ViewChild } from '@angular/core';

import { MainGroupService } from 'main-group/service/main-group.service';
import { MainGroup } from 'main-group/main-group';
import { MainGroupEditComponent } from 'main-group/main-group-edit.component';

@Component({
	selector: 'main-group',
	templateUrl: 'main-group.component.html',
	styleUrls: [ 'main-group.component.css' ],
	providers: [ MainGroupService ]
})

export class MainGroupComponent implements OnInit {
	mainGroups: MainGroup[];
	selectedMainGroup: MainGroup = null;
	errorMsg: string;

	@ViewChild(MainGroupEditComponent) editView: MainGroupEditComponent;

	constructor(private mainGroupService: MainGroupService) {

	}

	ngOnInit(): void {
		this.initMainGroups();
	}

	initMainGroups(): void {
		this.mainGroupService.getMainGroups().subscribe(
			data => this.mainGroups = data,
			error => this.errorMsg = <any> error
		);
	}

	selectMainGroup(mainGroup: MainGroup): void {
		this.selectedMainGroup = mainGroup;
	}

	onAction(): void {
		this.selectedMainGroup = null;
	}

	sumPercentage(): number {
		let result: number = 0;
		for (let group of this.mainGroups) {
			if (group._id != this.selectedMainGroup._id) {
				result += group.weight;
			}
		}
		return result;
	}
}

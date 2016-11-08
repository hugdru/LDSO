import { Component, OnInit } from '@angular/core';

import { MainGroupService } from 'main-group/service/main-group.service';
import { MainGroup } from 'main-group/main-group';

@Component({
	selector: 'main-group',
	templateUrl: 'main-group.component.html',
	styleUrls: ['main-group.component.css'],
	providers: [ MainGroupService ]
})

export class MainGroupComponent implements OnInit {
	mainGroups: MainGroup[];
	selectedMainGroup: MainGroup;
	errorMsg: string;

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

	updateMainGroup(): void {
		this.mainGroupService.updateMainGroup(this.selectedMainGroup);
		this.selectedMainGroup = null;
	}
}

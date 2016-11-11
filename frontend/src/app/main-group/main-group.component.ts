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
	selectedAddMainGroup: boolean = false;
	selectedShowSubGroup: MainGroup;
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

	selectAddMainGroup(): void {
		this.selectedAddMainGroup = true;
	}

	onShow(mainGroup: MainGroup): void {
		this.selectedShowSubGroup = mainGroup;
	}

	onDelete(mainGroup: MainGroup): void {
		this.mainGroupService.removeMainGroup(mainGroup._id).subscribe();
	}


}

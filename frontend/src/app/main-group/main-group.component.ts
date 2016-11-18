import { Component, OnInit, ViewChild } from '@angular/core';

import { MainGroupService } from 'main-group/service/main-group.service';
import { MainGroup } from 'main-group/main-group';
import { MainGroupEditComponent } from 'main-group/main-group-edit.component';

@Component({
	selector: 'main-group',
	templateUrl: 'html/main-group.component.html',
	styleUrls: [ 'main-group.component.css' ],
	providers: [ MainGroupService ]
})

export class MainGroupComponent implements OnInit {
	mainGroups: MainGroup[];
	parentMainGroup: MainGroup;
	errorMsg: string;

	constructor(private mainGroupService: MainGroupService) {

	}

	ngOnInit(): void {
		this.initMainGroups();
	}

	initMainGroups(): void {
		this.mainGroupService.getMainGroups().subscribe(
			data => this.mainGroups = data,
			error => this.errorMsg = <any>error
		);
	}

	onDelete(mainGroup: MainGroup): void {
		this.mainGroupService.removeMainGroup(mainGroup._id).subscribe();
	}

	onShow(mainGroup: MainGroup): void {
		this.parentMainGroup = mainGroup;
	}

}

import { Component, OnInit } from '@angular/core';

import { MainGroup } from 'main-group/main-group';
import { MainGroupService } from 'main-group/service/main-group.service';

@Component({
	selector: 'p4a-test-group',
	templateUrl: 'html/test-group.component.html',
	providers: [ MainGroupService ]
})

export class TestGroupComponent implements OnInit {
	groups: MainGroup[];
	errMsg: string;

	constructor(private mainGroupService: MainGroupService) { }

	ngOnInit() {
		this.initGroups();
	}

	initGroups(): void {
		this.mainGroupService.getMainGroups().subscribe(
			data => {this.groups = data; console.log(this.groups);},
			error => this.errMsg = <any>error
		);
	}
}

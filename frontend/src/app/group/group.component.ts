import { Component, OnInit } from '@angular/core';

import { Group } from 'group/group';
import { GroupService } from 'group/service/group.service';

@Component({
	selector: 'group-p4a',
	templateUrl: 'html/group.component.html'
	// styleUrls: [ 'group.component.css' ]
})

export class GroupComponent implements OnInit {
	groups: Group[];

	constructor(private groupService: GroupService ) { }

	ngOnInit(): void {
		this.initGroups();
	}

	setGroups(): void {
		this.GroupService.getGroup().subscribe(
			data => this.groups = data,
			error => this.errorMessage = <any>error
		);
	}
}

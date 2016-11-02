import { Component } from '@angular/core';

import { Group } from 'group/group';
import { GroupService } from 'group/service/group.service';

@Component({
	selector: 'group-p4a',
	templateUrl: 'html/group.component.html'
	// styleUrls: [ 'group.component.css' ]
})

export class GroupComponent {
	groups: Group[];

	constructor(private groupService: GroupService ) {}

	setGroups(): void {
		this.GroupService.getGroup().subscribe(data => this.groups = data);
	}
}

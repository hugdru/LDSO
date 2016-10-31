import { Component } from '@angular/core';

@Component({
	selector: 'group-p4a',
	templateUrl: 'html/group.component.html',
	styleUrls: [ 'group.component.css' ]
})

import { Group } from './group';
import { GroupService } from './service/group.service';

export class GroupComponent {
	groups: Group[];

	constructor(private groupService: GroupService ) {}

	setGroups(): void {
		this.GroupService.getGroup().subscribe(data => this.groups = data);
	}
}

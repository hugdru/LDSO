import { Component, OnInit, Input } from '@angular/core';

import { MainGroupService } from 'main-group/service/main-group.service';
import { MainGroup } from 'main-group/main-group';
import { SubGroup } from 'sub-group/sub-group';

@Component({
	selector: 'audit-list-groups',
	templateUrl: 'html/audit-list-groups.component.html',
	styleUrls: [ '../main-group/main-group.component.css' ],
	providers: [ ]
})

export class AuditListGroupsComponent implements OnInit {
	@Input() subGroups: SubGroup[];
	mainGroups: MainGroup[];
	mainGroupsId: number[];
	selectedMainGroup: MainGroup;

	constructor(private mainGroupService: MainGroupService) {
	}

	ngOnInit(): void {
		this.findMainGroups();
		this.getMainGroups();
	}

	getMainGroups(): void {
		for(let mainGroupId of this.mainGroupsId) {
			this.mainGroupService.getMainGroup("_id", "int", mainGroupId)
					.subscribe(data => this.mainGroups.push(data));
		}
	}

	findMainGroups(): void {
		for(let subGroup of this.subGroups) {
			if(!this.mainGroupsId.includes(subGroup.main_group)) {
				this.mainGroupsId.push(subGroup.main_group);
			}
		}
	}
}

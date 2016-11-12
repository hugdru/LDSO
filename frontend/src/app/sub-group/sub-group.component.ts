import {
	Component,
	OnInit,
	OnChanges,
	SimpleChanges,
	Input
} from '@angular/core';

import { SubGroupService } from 'sub-group/service/sub-group.service';
import { SubGroup } from 'sub-group/sub-group';
import { MainGroup } from 'main-group/main-group';

@Component({
	selector: 'sub-group',
	templateUrl: './html/sub-group.component.html',
	styleUrls: [ '../main-group/main-group.component.css' ],
	providers: [ SubGroupService ]
})

export class SubGroupComponent implements OnInit, OnChanges {
	subGroups: SubGroup[];
	selectedShowCriteria: SubGroup;

	@Input() selectedShowSubGroup: MainGroup;

	constructor(private subGroupService: SubGroupService){ }

	ngOnChanges(changes: SimpleChanges): void {
		for(let i in changes) {
			this.initSubGroups(changes[i].currentValue._id);
		}
	}

	ngOnInit() {
		this.initSubGroups(this.selectedShowSubGroup._id);
	}

	initSubGroups(mainGroupId: number): void {
		this.subGroupService
				.getSomeSubGroups("main_group", "int", mainGroupId)
				.subscribe(data => this.subGroups = data);

	}

	onDelete(subGroup: SubGroup): void {
		this.subGroupService.removeSubGroup(subGroup._id).subscribe();
	}

	onShow(subGroup: SubGroup): void {
		this.selectedShowCriteria = subGroup;
	}

}

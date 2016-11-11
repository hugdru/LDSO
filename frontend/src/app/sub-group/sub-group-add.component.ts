import { Component, Input, Output, EventEmitter, OnInit } from '@angular/core';

import { SubGroupService } from 'sub-group/service/sub-group.service';
import { SubGroup } from 'sub-group/sub-group';

@Component({
	selector: 'sub-group-add',
	templateUrl: '../main-group/html/main-group-edit.component.html',
	styleUrls: [ '../main-group/main-group-edit.component.css' ],
	providers: [ SubGroupService ]
})

export class SubGroupAddComponent implements OnInit {
	selectedObject: SubGroup;

	@Input() weight: number;
	@Output() onAdd = new EventEmitter<SubGroup>();

	constructor(private subGroupService: SubGroupService) {

	}

	ngOnInit(): void {
		this.selectedObject = new SubGroup();
	}

	pressed(newSubGroup: SubGroup): void {
		if(newSubGroup) {
			this.addSubGroup();
		}
		this.onAdd.emit(newSubGroup);
	}

	addSubGroup(): void {
		this.subGroupService.setSubGroup(this.selectedObject).subscribe();
	}

	checkPercentage(): boolean {
		return this.selectedObject.weight + this.weight > 100;
	}
}

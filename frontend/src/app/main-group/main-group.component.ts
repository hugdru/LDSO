import { Component, OnInit, ViewChild, AfterViewChecked } from '@angular/core';

import { MainGroupService } from 'main-group/service/main-group.service';
import { MainGroup } from 'main-group/main-group';
import { MainGroupEditComponent } from 'main-group/main-group-edit.component';

@Component({
	selector: 'main-group',
	templateUrl: 'main-group.component.html',
	styleUrls: ['main-group.component.css'],
	providers: [ MainGroupService ]
})

export class MainGroupComponent implements OnInit, AfterViewChecked {
	mainGroups: MainGroup[];
	selectedMainGroup: MainGroup = null;
	selectedAddMainGroup: boolean = false;
	errorMsg: string;

	@ViewChild(MainGroupEditComponent) editView: MainGroupEditComponent;

	constructor(private mainGroupService: MainGroupService) {

	}

	ngOnInit(): void {
		this.initMainGroups();
	}

	ngAfterViewChecked(): void {
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

	deleteMainGroup(mainGroup: MainGroup): void {
		this.mainGroupService.removeMainGroup(mainGroup._id).subscribe();
		let position: number;
		for(let i in this.mainGroups) {
			if(this.mainGroups[i]._id = mainGroup._id) {
				position = Number(i);
				break;
			}
		}
		this.mainGroups.splice(position, 1);
	}

	onAction(): void {
		this.selectedMainGroup = null;
	}

	onAdd(newMainGroup: MainGroup): void {
		if(newMainGroup != null) {
			this.mainGroups.push(newMainGroup);
		}
		this.selectedAddMainGroup = false;
	}

	checkPercentage(): boolean {
		let result: number;
		for (let group of this.mainGroups) {
			result += group.weight;
		}
		console.log(result <= 100);
		return result <= 100;
	}
}

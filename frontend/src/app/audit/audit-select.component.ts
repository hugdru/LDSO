import { Component, OnInit } from '@angular/core';
import { ActivatedRoute} from '@angular/router';

import { MainGroup } from 'main-group/main-group';
import { SubGroup } from "../sub-group/sub-group";
import { MainGroupService } from 'main-group/service/main-group.service';
import { SubGroupService } from "../sub-group/service/sub-group.service";

@Component({
	selector: 'p4a-audit',
	templateUrl: './html/audit-select.component.html',
	styleUrls: ['./audit.component.css'],
	providers: [MainGroupService, SubGroupService]
})

export class AuditSelectComponent implements OnInit {

	property_id: number;
	mainGroups: MainGroup[];
	subGroups: SubGroup[];
	errorMsg: string;
	selectedSubgroups: number[];

	constructor(
		private mainGroupService: MainGroupService,
		private subGroupService: SubGroupService,
		private route: ActivatedRoute
	) { }

	ngOnInit(): void {
		this.property_id = +this.route.snapshot.params['id'];
		this.initMainGroups();
        this.selectedSubgroups = [];
	}

	initMainGroups(): void {
		this.mainGroupService.getMainGroups().subscribe(
			data => this.mainGroups = data,
			error => this.errorMsg = <any> error
		);
	}

	showChildren(mainGroup: MainGroup): void {
		this.subGroupService.getSomeSubGroups("main_group", "int", mainGroup._id).subscribe(
			data => this.subGroups = data
		);
	}

    toggleSubGroup(id: number): void {
        var index = this.selectedSubgroups.indexOf(id);
        if (index > -1) {
            this.selectedSubgroups.splice(index, 1);
        }
        else {
            this.selectedSubgroups.push(id);
        }
    }

   checkedSubGroup(id: number): boolean {
       if (this.selectedSubgroups.indexOf(id) > -1) {
           return true;
       }
       return false;
   }

}

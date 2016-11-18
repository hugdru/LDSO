import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Params} from '@angular/router';

import { Property } from 'property/property';
import { PropertyService } from 'property/service/property.service';
import { MainGroup } from 'main-group/main-group';
import { MainGroupService } from 'main-group/service/main-group.service';

@Component({
	selector: 'p4a-audit',
	templateUrl: './html/audit.component.html',
	styleUrls: ['./audit.component.css'],
	providers: [ PropertyService ]
})

export class AuditComponent implements OnInit {

	property: Property;
	mainGroups: MainGroup[];
	errorMsg: string;

	constructor(
		private propertyService: PropertyService,
		private mainGroupService: MainGroupService,
		private route: ActivatedRoute
	) { }

	ngOnInit(): void {
		// this.route.params
		// 	.switchMap((params: Params) => this.propertyService.getProperty("_id", "int", +params['id']))
		// 	.subscribe(data => this.property = data);
		this.initMainGroups();
	}

	initMainGroups(): void {
		this.mainGroupService.getMainGroups().subscribe(
			data => this.mainGroups = data,
			error => this.errorMsg = <any> error
		);
	}

}

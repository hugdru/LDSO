import { Component, Input, OnInit } from '@angular/core';

import { Property } from './property';
import { PropertyService } from './service/property.service';

@Component({
	selector: 'properties-info',
	templateUrl: './html/properties-info.component.html',
	providers: [ PropertyService ]
})

export class PropertiesInfoComponent implements OnInit {
	property: Property;
	// @Input() property_id: number;
	property_id: number = 2;

	constructor(private propertyService: PropertyService) {
	}

	ngOnInit(): void{
		this.propertyService.getProperty("_id", "int", this.property_id)
				.subscribe(data => this.property = data);
	}

}

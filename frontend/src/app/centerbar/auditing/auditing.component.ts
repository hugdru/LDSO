import { Component, OnInit } from '@angular/core';

import { Criterion } from './criterion';
import { CriteriaService } from './criteria.service';

import { Property } from './property';
import { PropertyService } from './property.service';

const CRITERIA: Criterion[] = [
{ weight: 11, name: 'Mr. Nice' },
{ weight: 12, name: 'Narco' },
{ weight: 13, name: 'Bombasto' },
{ weight: 14, name: 'Celeritas' },
{ weight: 15, name: 'Magneta' },
{ weight: 16, name: 'RubberMan' },
{ weight: 17, name: 'Dynama' },
{ weight: 18, name: 'Dr IQ' },
{ weight: 19, name: 'Magma' },
{ weight: 20, name: 'Tornado' }
];


@Component({
  selector: 'p4a-auditing',
  templateUrl: './auditing.component.html',
  styleUrls: ['./auditing.component.css'],
	providers: [ CriteriaService , PropertyService ]
})
export class AuditingComponent implements OnInit {


	constructor() { }

	ngOnInit(): void {

	}



}
export class PropertiesInfoComponent implements OnInit {
	property: Property;

	constructor(private propertyService: PropertyService) { }

	ngOnInit(): void {
		this.getProperty();
	}

	getProperty(): void {
		this.propertyService.getProperty('Hotel Sunny')
				.subscribe(data => this.property = data);
	}
}
export class CriteriaSelectorComponent implements OnInit {
	criteria : Criterion[];
	selectedCriterion: Criterion;

	constructor(private criteriaService: CriteriaService) { }

	ngOnInit(): void {
		this.getCriteria();
	}

	getCriteria(): void {
		this.criteriaService.getCriteria()
				.subscribe(data => this.criteria = data);
	}

	showSubCriteria(criterion: Criterion): void {
		this.selectedCriterion = criterion;
	}

}




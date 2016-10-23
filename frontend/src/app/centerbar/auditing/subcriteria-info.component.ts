import { Component, Input, OnInit } from '@angular/core';

import { Criterion } from './criterion';
// import { CriteriaService } from './criterion.service';

@Component({
	selector: 'subcriteria-info',
	templateUrl: 'app/audit/html/subcriteria-info.component.html',
	// styleUrls: [ './css/subcriteria-info.component.css' ],
	// providers: [ CriteriaService ]
})

export class SubCriteriaInfoComponent {
	@Input() selectedCriterion: Criterion;

}

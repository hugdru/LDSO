import { Component, OnInit } from '@angular/core';

import { SubGroupService } from 'sub-group/service/sub-group.service';
import { SubGroup } from 'sub-group/sub-group';

@Component({
	selector: 'app-sub-group',
	templateUrl: './html/sub-group.component.html',
	styleUrls: [ 'sub-group.component.css' ],
	providers: [ SubGroupService ]
})

export class SubGroupComponent implements OnInit {

	ngOnInit() {

	}

}

import { Component, Input, } from '@angular/core';

import { Identifier } from 'identifier.interface';

@Component({
	selector: 'list-show-value',
	templateUrl: './list-show-value.component.html',
	styleUrls: [ '../../main-group/main-group.component.css' ],
})

export class ListShowValueComponent {
	@Input() object: Identifier;
	@Input() objType: string;
}


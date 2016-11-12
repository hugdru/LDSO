import { Component, Input, Output, EventEmitter } from '@angular/core';

import { Criterion } from 'criterion/criterion';
import { SubGroup } from 'sub-group/sub-group';
import { MainGroup } from 'main-group/main-group';
import { Identifier } from 'identifier.interface';

@Component({
	selector: 'list-select-add',
	templateUrl: './list-select-add.component.html',
	styleUrls: [ '../../main-group/main-group.component.css' ],
})

export class ListSelectAddComponent {
	@Input() objType: string;
	@Input() weight: number;
	@Input() father: Identifier;
	@Output() add = new EventEmitter<Identifier>();

	onAdd(newObject: Identifier): void {
		this.add.emit(newObject);
	}
}


import { Component, Input, Output, EventEmitter } from '@angular/core';

import { Criterion } from 'criterion/criterion';
import { SubGroup } from 'sub-group/sub-group';
import { MainGroup } from 'main-group/main-group';
import { Identifier } from 'identifier.interface';

@Component({
	selector: 'list-select-edit',
	templateUrl: './list-select-edit.component.html',
	styleUrls: [ '../../main-group/main-group.component.css' ],
})

export class ListSelectEditComponent {
	@Input() objType: string;
	@Input() selectedEditObject: Identifier;
	@Input() weight: number;
	@Output() action = new EventEmitter();

	onAction(): void {
		this.action.emit();
	}
}


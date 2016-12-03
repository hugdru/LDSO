import {Component, Input, Output, EventEmitter} from "@angular/core";

@Component({
    selector: 'show-list',
    templateUrl: './html/show-list.component.html',
    // styleUrls: ['./audit.component.css'],
    // providers: [MainGroupService, SubGroupService]
})

export class ShowListComponent {
	
	@Input() objectList: Object[];
	@Output() selected = new EventEmitter<Object>();

	onSelect(selectedObject: Object): void {
		this.selected.emit(selectedObject);
	}
}

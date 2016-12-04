import {Component, Input, Output, EventEmitter} from "@angular/core";
import {Identifier} from "identifier.interface";

@Component({
    selector: 'list-select-edit',
    templateUrl: './list-select-edit.component.html',
    styleUrls: ['../../main-group/main-group.component.css'],
})

export class ListSelectEditComponent {
    @Input() objType: string;
    @Input() selectedEditObject: Identifier;
    @Input() weight: number;
    @Input() father: Identifier;
    @Output() action = new EventEmitter();

    onAction(): void {
        this.action.emit();
    }
}


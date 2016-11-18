import {
    Component,
    Input,
    Output,
    EventEmitter
} from '@angular/core';

import { SubGroup } from 'sub-group/sub-group';
import { MainGroup } from 'main-group/main-group';
import { Identifier } from 'identifier.interface';

@Component({
    selector: 'list-select-audit',
    templateUrl: './list-select-audit.component.html',
    styleUrls: [ '../../main-group/main-group.component.css' ],
})

export class ListSelectAuditComponent {
    selectedEditObject: Identifier;

    @Input() objects: Identifier[];
    @Input() father: Identifier;
    @Output() onShow = new EventEmitter<Object>();
    @Output() onDelete = new EventEmitter<Object>();

    showChildren(obj: Identifier){
        this.onShow.emit(obj);
    }

    findType(): string {
        if(this.father === undefined) {
            return "MainGroup";
        }
        else {
            return "SubGroup";
        }
    }

    onAction(): void {
        this.selectedEditObject = null;
    }

}

import {Component, Input, Output, EventEmitter} from "@angular/core";
import {Identifier} from "identifier.interface";

@Component({
    selector: 'list-manage',
    templateUrl: './list-manage.component.html',
    styleUrls: ['../../audit-template/audit-template.component.css'],
})

export class ListManageComponent {
    selectedEditObject: Identifier;
    selectedAddObject: boolean = false;

    @Input() objType: string;
    @Input() objects: Identifier[];
    @Input() father: Identifier;
    @Output() onShow = new EventEmitter<Object>();
    @Output() onDelete = new EventEmitter<Object>();
    @Output() onClose = new EventEmitter<Object>();
    @Output() onOpen = new EventEmitter<Object>();

    showChildren(obj: Identifier) {
        this.onShow.emit(obj);
    }

    deleteObject(obj: Identifier): void {
        this.onDelete.emit(obj);
        let position: number;
        for (let i in this.objects) {
            if (this.objects[i].id == obj.id) {
                position = Number(i);
                break;
            }
        }
        this.objects.splice(position, 1);
    }

    closeObject(obj: Identifier): void {
        this.onClose.emit(obj);
    }

    openObject(obj: Identifier): void {
        this.onOpen.emit(obj);
    }

    selectAddObject() {
        this.selectedAddObject = true;
    }

    selectEditObject(obj: Identifier) {
        this.selectedEditObject = obj;
    }

    sumPercentageForAdd(): number {
        let result: number = 0;
        for (let obj of this.objects) {
            result += obj.weight;
        }
        return result;
    }

    sumPercentage(): number {
        let result: number = 0;
        for (let obj of this.objects) {
            if (obj.id != this.selectedEditObject.id) {
                result += obj.weight;
            }
        }
        return result;
    }

    checkPercentage(): boolean {
        return this.sumPercentageForAdd() != 100;
    }

    onAction(): void {
        this.selectedEditObject = null;
    }

    onAdd(newObject: Identifier): void {
        if (newObject != null) {
            this.objects.push(newObject);
        }
        this.selectedAddObject = false;
    }

}

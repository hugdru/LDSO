import {Component, Input, Output, EventEmitter} from "@angular/core";
import {Criterion} from "criterion/criterion";
import {SubGroup} from "sub-group/sub-group";
import {MainGroup} from "main-group/main-group";
import {Identifier} from "identifier.interface";

@Component({
    selector: 'list-manage',
    templateUrl: './list-manage.component.html',
    styleUrls: ['../../ctemplate/ctemplate.component.css'],
})

export class ListManageComponent {
    selectedEditObject: Identifier;
    selectedAddObject: boolean = false;

    @Input() objects: Identifier[];
    @Input() father: Identifier;
    @Output() onShow = new EventEmitter<Object>();
    @Output() onDelete = new EventEmitter<Object>();

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
        return this.sumPercentageForAdd() > 100;
    }

    findType(): string {
        if (this.father === undefined) {
            return "Ctemplate";
        }
        else if ((<MainGroup>this.father).idTemplate !== undefined) {
            return "SubGroup";
        }
        else if ((<SubGroup>this.father).idMaingroup !== undefined) {
            return "Criterion";
        }
        else if ((<Criterion>this.father).idSubgroup !== undefined) {
            return "Accessibility";
        }
        else {
            return "MainGroup";
        }
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

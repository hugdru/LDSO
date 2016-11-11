import {
	Component,
	OnInit,
	Input,
	Output,
	EventEmitter
} from '@angular/core';

import { SubGroup } from 'sub-group/sub-group';
import { MainGroup } from 'main-group/main-group';

@Component({
	selector: 'list-manage',
	templateUrl: './list-manage.component.html',
	styleUrls: [ '../../main-group/main-group.component.css' ],
})

export class ListManageComponent implements OnInit {
	selectedEditObject: Object;
	isMainGroup: boolean = false;
	isSubGroup: boolean = false;
	selectedAddObject: boolean = false;

	@Input() objects: Object[];
	@Output() onShow = new EventEmitter<Object>();
	@Output() onDelete = new EventEmitter<Object>();

	ngOnInit(): void {
		if(this.objects instanceof MainGroup) {
			this.isMainGroup = true;
		} else {
			this.isSubGroup = true;
		}
	}

	showChildren(obj: Object){
		this.onShow.emit(obj);
	}

	deleteObject(obj: Object): void{
		this.onDelete.emit(obj);
		let position: number;
		if(this.isMainGroup) {
			let objects = this.objects as MainGroup[];
			let object = obj as MainGroup;
			for(let i in objects) {
				if(objects[i]._id = object._id) {
					position = Number(i);
					break;
				}
			}
		} else {
			let objects = this.objects as SubGroup[];
			let object = obj as SubGroup;
			for(let i in objects) {
				if(objects[i]._id = object._id) {
					position = Number(i);
					break;
				}
			}
		}
		this.objects.splice(position, 1);
	}

	onAdd(newObject: Object): void {
		if(newObject != null) {
			this.objects.push(newObject);
		}
		this.selectedAddObject = false;
	}

	selectAddObject(){
		this.selectedAddObject = true;
	}

	selectEditObject(obj: Object){
		this.selectedEditObject = obj;
	}

	sumPercentageForAdd(): number {
		let result: number = 0;

		if(this.isMainGroup) {
			let objects = this.objects as MainGroup[];
			for (let obj of objects) {
				result += obj.weight;
			}
		} else {
			let objects = this.objects as SubGroup[];
			for (let obj of objects) {
				result += obj.weight;
			}
		}
		return result;
	}

	sumPercentage(): number {
		let result: number = 0;
		if(this.isMainGroup) {
			let objects = this.objects as MainGroup[];
			let object = this.selectedEditObject as MainGroup;
			for (let obj of objects) {
				if (obj._id != object._id) {
					result += obj.weight;
				}
			}
		} else {
			let objects = this.objects as SubGroup[];
			let object = this.selectedEditObject as SubGroup;
			for (let obj of objects) {
				if (obj._id != object._id) {
					result += obj.weight;
				}
			}
		}
		return result;
	}

	onAction(): void {
		this.selectedEditObject = null;
	}

}

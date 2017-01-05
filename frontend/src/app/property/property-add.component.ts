import {Component, Output, EventEmitter, OnInit} from "@angular/core";
import {PropertyService} from "./service/property.service";
import {Property} from "./property";

@Component({
    selector: 'property-add',
    templateUrl: '/html/property-add.component.html',
    providers: [PropertyService]
})

export class PropertyAddComponent implements OnInit {
    selectedObject: Property;

    constructor(private propertyService: PropertyService) {

    }

    ngOnInit(): void {
        this.selectedObject = new Property();
    }

    addProperty(): void {
        this.propertyService.setProperty(this.selectedObject).subscribe(
                response => this.selectedObject.id = response.json().id
        );
    }

    cancel(): void {
        this.selectedObject = new Property();
    }

}

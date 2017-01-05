import {Component, OnInit} from "@angular/core";
import {PropertyService} from "./service/property.service";
import {Property} from "./property";
import {Address} from "./address";
import {Router} from "@angular/router";

@Component({
    selector: 'property-add',
    templateUrl: '/html/property-add.component.html',
    providers: [PropertyService]
})

export class PropertyAddComponent implements OnInit {
    selectedObject: Property;

    constructor(private propertyService: PropertyService, private router: Router) {

    }

    ngOnInit(): void {
        this.selectedObject = new Property();
        this.selectedObject.address = new Address();
    }

    addProperty(): void {
        this.propertyService.setProperty(this.selectedObject).subscribe(
                response => {
                    this.selectedObject.id = response.json().id;
                    this.router.navigate(['/listPropertiesView']);
                }
        );
    }

    cancel(): void {
        this.selectedObject = new Property();
        this.selectedObject.address = new Address();
    }

}
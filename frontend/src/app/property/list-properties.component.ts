import {Component, OnInit} from "@angular/core";
import {Property} from "property/property";
import {PropertyService} from "property/service/property.service";
import {Router} from "@angular/router";

@Component({
    selector: 'list-properties',
    templateUrl: 'html/list-properties.component.html',
    providers: [PropertyService]
})

export class ListPropertiesComponent implements OnInit {

    properties: Property[];
    router: Router;

    constructor(private propertyService: PropertyService, private _router: Router) {
        this.router = _router;
    }

    ngOnInit(): void {
        this.initProperties();
    }

    initProperties(): void {
        this.propertyService.getProperties().subscribe(
                data => this.properties = data
        );
    }

}


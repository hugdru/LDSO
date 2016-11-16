import {
    Component,
    OnInit
} from '@angular/core';

import { Property } from 'property/property';
import { PropertyService } from 'property/service/property.service';

@Component({
    selector: 'list-properties',
    templateUrl: 'html/list-properties.component.html',
    providers: [ PropertyService ]
})

export class ListPropertiesComponent implements OnInit {

    properties: Property[];

    constructor(private propertyService: PropertyService) {

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


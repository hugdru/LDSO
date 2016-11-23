import {Observable} from "rxjs/Rx";
import {ListPropertiesComponent} from "property/list-properties.component";
import {PropertyService} from "property/service/property.service";

class MockPropertyService extends PropertyService {
    constructor() {
        super(null);
    }

    getProperties() {
        return Observable.of([
            {_id: 26, name: "ana", image_path: "bla"},
            {_id: 14, name: "joao", image_path: "ble"}
        ]);
    }
}

describe('Property unit test', () => {
    let listPropertiesComponent: ListPropertiesComponent,
            propertyService: PropertyService;

    beforeEach(() => {
        propertyService = new MockPropertyService();
        listPropertiesComponent = new ListPropertiesComponent(propertyService);
    });

    it('shows list of property items by default - unit', () => {
        listPropertiesComponent.ngOnInit();
        expect(listPropertiesComponent.properties.length).toBe(2);
    });
});
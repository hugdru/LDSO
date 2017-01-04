import {Response} from "@angular/http";
import {Injectable} from "@angular/core";
import {Observable} from "rxjs/Observable";
import {propertiesUrl} from "shared/shared-data";
import {HandlerService} from "../../shared/service/handler.service";
import {Property} from "property/property";

@Injectable()
export class PropertyService {

    constructor(private handler: HandlerService) {
    }

    getProperties(): Observable<Property[]> {
        return this.handler.getAll<Property[]>(propertiesUrl);
    }

    getSomeProperties(tag: string, value: any): Observable<Property[]> {
        return this.handler.getSome<Property[]>(propertiesUrl, tag, value);
    }

    getProperty(id: number): Observable<Property> {
        return this.handler.get<Property>(propertiesUrl, id);
    }

    updateProperty(property: Property): Observable<Response> {
        return this.handler.update<Property>(propertiesUrl, property, property.id);
    }

    setProperty(property: Property): Observable<Response> {
        return this.handler.set<Property>(propertiesUrl, property);
    }

    removeProperty(id: number): Observable<Response> {
        return this.handler.delete(propertiesUrl, id);
    }
}

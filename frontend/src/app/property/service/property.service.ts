import {Response} from "@angular/http";
import {Injectable} from "@angular/core";
import {Observable} from "rxjs/Observable";
import {propertiesUrl, propertiesFindUrl} from "shared/shared-data";
import {HandlerService} from "handler.service";
import {Property} from "property/property";

@Injectable()
export class PropertyService {

    constructor(private handler: HandlerService) {
    }

    getProperties(): Observable<Property[]> {
        return this.handler.getAll<Property[]>(propertiesUrl);
    }

    getSomeProperties(tag: string, type: string, value: any): Observable<Property[]> {
        return this.handler.get<Property[]>(propertiesUrl, tag, type,
                value);
    }

    getProperty(tag: string, type: string, value: any): Observable<Property> {
        return this.handler.get<Property>(propertiesFindUrl, tag, type, value);
    }

    updateProperty(property: Property): Observable<Response> {
        return this.handler.update<Property>(propertiesUrl, property, property._id);
    }

    setProperty(property: Property): Observable<Response> {
        return this.handler.set<Property>(propertiesUrl, property);
    }

    removeProperty(id: number): Observable<Response> {
        return this.handler.delete(propertiesUrl, id);
    }
}

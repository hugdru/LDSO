import {Response} from "@angular/http";
import {Injectable} from "@angular/core";
import {Observable} from "rxjs/Observable";
import {accessibilitiesUrl, accessibilitiesFindUrl} from "shared/shared-data";
import {HandlerService} from "handler.service";
import {Accessibility} from "accessibility/accessibility";

@Injectable()
export class AccessibilityService {

    constructor(private handler: HandlerService) {
    }

    getAccessibilities(): Observable<Accessibility[]> {
        return this.handler.getAll<Accessibility[]>(accessibilitiesUrl);
    }

    getSomeAccessibilities(tag: string, type: string, value: any): Observable<Accessibility[]> {
        return this.handler.get<Accessibility[]>(accessibilitiesUrl,
                tag, type, value);
    }

    getAccessibility(tag: string, type: string, value: any): Observable<Accessibility> {
        return this.handler.get<Accessibility>(accessibilitiesFindUrl, tag, type,
                value);
    }

    updateAccessibility(accessibility: Accessibility): Observable<Response> {
        return this.handler.update<Accessibility>(accessibilitiesUrl,
                accessibility, accessibility._id);
    }

    setAccessibility(accessibility: Accessibility): Observable<Response> {
        return this.handler.set<Accessibility>(accessibilitiesUrl,
                accessibility);
    }

    removeAccessibility(id: number): Observable<Response> {
        return this.handler.delete(accessibilitiesUrl, id);
    }

}

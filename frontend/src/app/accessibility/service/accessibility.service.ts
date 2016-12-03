import {Response} from "@angular/http";
import {Injectable} from "@angular/core";
import {Observable} from "rxjs/Observable";
import {accessibilitiesUrl} from "shared/shared-data";
import {HandlerService} from "handler.service";
import {Accessibility} from "accessibility/accessibility";

@Injectable()
export class AccessibilityService {

    constructor(private handler: HandlerService) {
    }

    getAccessibilities(): Observable<Accessibility[]> {
        return this.handler.getAll<Accessibility[]>(accessibilitiesUrl);
    }

    getSomeAccessibilities(id: number): Observable<Accessibility[]> {
        console.log(accessibilitiesUrl.replace(/#/g, id.toString()));
        return this.handler.getAll<Accessibility[]>(
                accessibilitiesUrl.replace(/#/g, id.toString()));
    }

    getAccessibility(id: number): Observable<Accessibility> {
        return this.handler.get<Accessibility>(accessibilitiesUrl, id);
    }

    updateAccessibility(accessibility: Accessibility): Observable<Response> {
        return this.handler.update<Accessibility>(accessibilitiesUrl,
                accessibility, accessibility.id);
    }

    setAccessibility(accessibility: Accessibility): Observable<Response> {
        return this.handler.set<Accessibility>(accessibilitiesUrl,
                accessibility);
    }

    removeAccessibility(id: number): Observable<Response> {
        return this.handler.delete(accessibilitiesUrl, id);
    }

}

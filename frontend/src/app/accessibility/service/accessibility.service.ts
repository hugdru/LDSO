import {Response} from "@angular/http";
import {Injectable} from "@angular/core";
import {Observable} from "rxjs/Observable";
import {accessibilitiesUrl} from "shared/shared-data";
import {HandlerService} from "../../shared/service/handler.service";
import {Accessibility} from "accessibility/accessibility";

@Injectable()
export class AccessibilityService {

    constructor(private handler: HandlerService) {
    }

    getAccessibilities(): Observable<Accessibility[]> {
        return this.handler.getAll<Accessibility[]>(accessibilitiesUrl);
    }

    getSomeAccessibilities(criterionId: number): Observable<Accessibility[]> {
        return this.handler.getAll<Accessibility[]>(
                accessibilitiesUrl.replace(/#/g, criterionId.toString()));
    }

    getAccessibility(id: number): Observable<Accessibility> {
        return this.handler.get<Accessibility>(accessibilitiesUrl, id);
    }

    updateAccessibility(accessibility: Accessibility,
                        criterionId: number): Observable<Response> {
        return this.handler.update<Accessibility>(
                accessibilitiesUrl.replace(/#/g, criterionId.toString()),
                accessibility, accessibility.id);
    }

    setAccessibility(accessibility: Accessibility,
                     criterionId: number): Observable<Response> {
        return this.handler.set<Accessibility>(
                accessibilitiesUrl.replace(/#/g, criterionId.toString()),
                accessibility);
    }

    removeAccessibility(id: number, criterionId: number): Observable<Response> {
        return this.handler.delete(
                accessibilitiesUrl.replace(/#/g, criterionId.toString()), id);
    }

}

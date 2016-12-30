import {Response} from "@angular/http";
import {Injectable} from "@angular/core";
import {Observable} from "rxjs/Observable";
import {loginUrl, logoutUrl} from "shared/shared-data";
import {HandlerService} from "handler.service";
import {Session} from "../session";

@Injectable()
export class LoginService {
    constructor(private handler: HandlerService) {
    }

    setLogin(session: Session): Observable<Response> {
        return this.handler.set<Session>(loginUrl, session);
    }

    getLogout(): Observable<Response> {
        return this.handler.getResponse(logoutUrl);
    }

}

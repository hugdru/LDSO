import {Response} from "@angular/http";
import {Injectable} from "@angular/core";
import {Observable} from "rxjs/Observable";
import {loginUrl} from "shared/shared-data";
import {HandlerService} from "handler.service";
import {Login} from "login/login";

@Injectable()
export class LoginService {
    constructor(private handler: HandlerService) {
    }

    setLogin(login: Login): Observable<Response> {
        return this.handler.set<Login>(loginUrl, login);
    }
}

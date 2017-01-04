import {Response} from "@angular/http";
import {Injectable} from "@angular/core";
import {Observable} from "rxjs/Observable";
import {registerUrl} from "shared/shared-data";
import {HandlerService} from "../../shared/service/handler.service";
import {User} from "../user";

@Injectable()
export class RegisterService {
    constructor(private handler: HandlerService) {
    }

    setRegister(user: User): Observable<Response> {
        return this.handler.set<User>(registerUrl, user);
    }

}

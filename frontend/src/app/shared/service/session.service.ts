import {Injectable} from '@angular/core';
import {Session} from "../../login/session";

@Injectable()
export class SessionService {

    private loggedIn: boolean;
    private session: Session;

    constructor() {
        this.loggedIn = !!localStorage.getItem('auth_token');
        this.initSession(this.loggedIn)
    }

    initSession(loggedIn: boolean): void {
        if (loggedIn)
            this.session = JSON.parse(localStorage.getItem('session'));
        else {
            this.session = new Session();
        }
    }

    getSession(): Session {
        this.loggedIn = !!localStorage.getItem('auth_token');
        this.initSession(this.loggedIn);
        return this.session;
    }

}



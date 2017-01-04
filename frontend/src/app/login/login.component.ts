import {Component, OnInit} from "@angular/core";
import {LoginService} from "login/service/login.service";
import {Session} from "./session";
import {SessionService} from "../shared/service/session.service";

@Component({
    selector: 'login',
    templateUrl: './html/login.component.html',
    styleUrls: [],
    providers: [LoginService]

})

export class LoginComponent implements OnInit {

    session: Session;
    loggedIn = false;
    errorMsg: string;

    constructor(private loginService: LoginService,
                private sessionService: SessionService) {
        this.loggedIn = !!localStorage.getItem('auth_token');
    }

    ngOnInit(): void {
        if (this.loggedIn)
            this.session = JSON.parse(localStorage.getItem('session'));
        else this.session = new Session();
    }

    login(): void {
        this.loginService.setLogin(this.session).subscribe(
                response => {
                    this.session.id = response.json().id;
                    this.session.role = response.json().role;
                    this.session.name = response.json().name;
                    this.session.email = response.json().email;
                    this.session.password = "";
                    localStorage.setItem('auth_token', response.json().auth_token);
                    localStorage.setItem('session', JSON.stringify(this.session));
                    this.loggedIn = true;
                    this.sessionService.announceSession(this.loggedIn);
                },
                error => {
                    this.errorMsg = <any>error;
                }
        );
    }

    logout(): void {
        this.loginService.getLogout().subscribe(
                error => {
                    this.errorMsg = <any>error;
                }
        );
        localStorage.removeItem('auth_token');
        localStorage.removeItem('session');
        this.session = new Session();
        this.loggedIn = false;
        this.sessionService.announceSession(this.loggedIn);
    }

    cancel(): void {
        this.session = new Session();
    }

}
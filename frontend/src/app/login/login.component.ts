import {Component, OnInit} from "@angular/core";
import {LoginService} from "login/service/login.service";
import {Session} from "./session";

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

    constructor(private loginService: LoginService) {
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
                    this.loggedIn = true;
                    localStorage.setItem('auth_token', response.json().auth_token);
                    localStorage.setItem('session', JSON.stringify(this.session));
                },
                error => {
                    this.errorMsg = <any>error;
                    console.log(this.errorMsg);
                }
        );
    }

    logout(): void {
        this.loginService.getLogout().subscribe(
                error => {
                    this.errorMsg = <any>error;
                    console.log(this.errorMsg);
                }

        );
        localStorage.removeItem('auth_token');
        localStorage.removeItem('session');
        this.session = new Session();
        this.loggedIn = false;
    }

}
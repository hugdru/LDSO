import {Component, Output, EventEmitter, OnInit} from "@angular/core";
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
    errorMsg: string;
//    @Output() onAdd = new EventEmitter<Session>();

    constructor(private loginService: LoginService) {
    }

    ngOnInit(): void {
        this.session = new Session();
    }

    loginPressed(newSession: Session): void {
        if (newSession) {
            this.login();
        }
//        this.onAdd.emit(newSession);
    }

    login(): void {
        this.loginService.setLogin(this.session).subscribe(
                response => {
                    this.session.id = response.json().id;
                    this.session.role = response.json().role;
                    this.session.name = response.json().name;
                    this.session.email = response.json().email;
                },
                error => {
                    this.errorMsg = <any>error;
                    console.log(this.errorMsg);
                }
        );
        this.session.password = "";
        console.log(this.session);
    }

    logoutPressed(newSession: Session): void {
        if (newSession) {
            this.logout();
        }
//        this.onAdd.emit(newSession);
    }

    logout(): void {
        this.loginService.getLogout().subscribe(
                response => this.session = new Session(),
                error => {
                    this.errorMsg = <any>error;
                    console.log(this.errorMsg);
                }
        );
        console.log(this.session);
    }

}
import {Component, Output, EventEmitter, OnInit} from "@angular/core";
import {LoginService} from "login/service/login.service";
import {Login} from "login/login";

@Component({
    selector: 'login',
    templateUrl: './html/login.component.html',
    styleUrls: [],
    providers: [LoginService]

})

export class LoginComponent implements OnInit {

    selectedObject: Login;
    errorMsg: string;
//    @Output() onAdd = new EventEmitter<Login>();

    constructor(private loginService: LoginService) {
    }

    ngOnInit(): void {
        this.selectedObject = new Login();
    }

    pressed(newLogin: Login): void {
        if (newLogin) {
            this.login();
        }
//        this.onAdd.emit(newLogin);
    }

    login(): void {
        this.loginService.setLogin(this.selectedObject).subscribe(
                response => {
                    this.selectedObject.id = response.json().id;
                    console.log(response.json());
                    console.log(this.selectedObject.id);
                },
                error => {
                    this.errorMsg = <any>error;
                    console.log(this.errorMsg);
                }
        );
    }

}
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
                    this.selectedObject.role = response.json().role;
                    this.selectedObject.name = response.json().name;
                    this.selectedObject.email = response.json().email;
                },
                error => {
                    this.errorMsg = <any>error;
                    console.log(this.errorMsg);
                }
        );
        this.selectedObject.password = "";
        console.log(this.selectedObject);
    }

}
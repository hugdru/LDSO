import {Component} from "@angular/core";
import {LoginService} from "login/service/login.service";
import {Login} from "login/login";

@Component({
    selector: 'login',
    templateUrl: './html/login.component.html',
    styleUrls: [],
    providers: [LoginService]

})

export class LoginComponent {

    selectedObject: Login;

    constructor(private loginService: LoginService) {
    }


    login(): void {
        this.loginService.setLogin(this.selectedObject).subscribe(
                response => this.selectedObject.id = response.json().id
        );
    }

}
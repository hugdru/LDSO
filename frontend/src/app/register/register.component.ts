import {Component, OnInit} from "@angular/core";
import {RegisterService} from "register/service/register.service";
import {User} from "./user";
import {Router} from "@angular/router";
import {Session} from "../login/session";

@Component({
    selector: 'register',
    templateUrl: './html/register.component.html',
    styleUrls: [],
    providers: [RegisterService]

})

export class RegisterComponent implements OnInit {

    user: User;
    session: Session;
    errorMsg: string;
    optionRolesSA = [
        {value: "client", label: "Cliente"},
        {value: "superadmin", label: "Administrador"},
        {value: "localadmin", label: "Administrador Local"},
        {value: "auditor", label: "Auditor"},
    ];
    optionRolesLA = [
        {value: "client", label: "Cliente"},
        {value: "auditor", label: "Auditor"},
    ];


    constructor(private registerService: RegisterService, private router: Router) {
    }

    ngOnInit(): void {
        this.user = new User();
        if (!!localStorage.getItem('auth_token'))
            this.session = JSON.parse(localStorage.getItem('session'));
        else this.session = new Session();
    }

    register(): void {
        if (!this.user.role) this.user.role = "client";
        this.registerService.setRegister(this.user).subscribe(
                response => {
                    this.user.id = response.json().id;
                    this.router.navigate(['/login']);
                },
                error => {
                    this.errorMsg = <any>error;
                }
        );
    }

    cancel(): void {
        this.user = new User();
    }

}
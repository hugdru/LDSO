import {Component, OnInit} from "@angular/core";
import {RegisterService} from "register/service/register.service";
import {User} from "./user";

@Component({
    selector: 'register',
    templateUrl: './html/register.component.html',
    styleUrls: [],
    providers: [RegisterService]

})

export class RegisterComponent implements OnInit {

    user: User;
    errorMsg: string;

    constructor(private registerService: RegisterService) {
    }

    ngOnInit(): void {
        this.user = new User();
    }

    register(): void {
        this.registerService.setRegister(this.user).subscribe(
                response => {
                    this.user.id = response.json().id;
                },
                error => {
                    this.errorMsg = <any>error;
                }
        );
        console.log(this.user);
    }

}
import {Injectable} from "@angular/core";
import {CanActivate} from "@angular/router";
import {SessionService} from "./session.service";
import {Session} from "../../login/session";

@Injectable()
export class AuditGuard implements CanActivate {

    private session: Session;

    constructor(private sessionService: SessionService) {}

    canActivate() {
        this.session = this.sessionService.getSession();
        console.log(this.session);
        return this.session != undefined && (
                this.session.role == 'superadmin' ||
                this.session.role == 'localadmin' ||
                this.session.role == 'auditor'
                );
    }

}
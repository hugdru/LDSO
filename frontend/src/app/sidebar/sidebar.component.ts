import {Component, OnInit} from "@angular/core";
import {Session} from "../login/session";
import {SessionService} from "../shared/service/session.service";
import {Subscription} from "rxjs";


@Component({
    selector: 'p4a-sidebar',
    templateUrl: './sidebar.component.html',
    styleUrls: ['./sidebar.component.css']

})

export class SidebarComponent implements OnInit {

    loggedIn: boolean;
    session: Session;
    subscription: Subscription;

    constructor(sessionService: SessionService) {
        this.loggedIn = !!localStorage.getItem('auth_token');
        this.subscription = sessionService.sessionAnnounced$.subscribe(
                loggedIn => {
                    console.log(loggedIn);
                    this.loggedIn = loggedIn;
                    this.initSession(this.loggedIn);
                });
    }

    ngOnInit(): void {
        this.initSession(this.loggedIn);
    }

    ngOnChanges(): void {
        this.initSession(this.loggedIn);
    }

    initSession(loggedIn: boolean): void {
        if (loggedIn)
            this.session = JSON.parse(localStorage.getItem('session'));
        else this.session = new Session();
    }

    setLogin(value: boolean): void {
        this.loggedIn = value;
    }

    ngOnDestroy() {
        this.subscription.unsubscribe();
    }

}


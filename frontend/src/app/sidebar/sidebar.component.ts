import {Component, OnInit} from "@angular/core";
import {Session} from "../login/session";
import {SessionAnnounceService} from "../shared/service/session-announce.service";
import {Subscription} from "rxjs";


@Component({
    selector: 'p4a-sidebar',
    templateUrl: './sidebar.component.html',
    styleUrls: ['./sidebar.component.css'],
})

export class SidebarComponent implements OnInit {

    loggedIn: boolean;
    session: Session;
    subscription: Subscription;

    constructor(private sessionService: SessionAnnounceService) {
        this.loggedIn = !!localStorage.getItem('auth_token');
        this.subscription = sessionService.sessionAnnounced$.subscribe(
                loggedIn => {
                    this.loggedIn = loggedIn;
                    this.initSession(loggedIn);
                });
    }

    ngOnInit(): void {
        this.initSession(this.loggedIn)
    }

    ngOnChanges(): void {
        this.initSession(this.loggedIn);
    }

    initSession(loggedIn: boolean): void {
        if (loggedIn)
            this.session = JSON.parse(localStorage.getItem('session'));
        else this.session = new Session();
    }

    ngOnDestroy() {
        this.subscription.unsubscribe();
    }

}


import {Injectable} from "@angular/core";
import {Subject} from "rxjs/Subject";

@Injectable()
export class SessionAnnounceService {

    private sessionAnnouncedSource = new Subject<boolean>();

    sessionAnnounced$ = this.sessionAnnouncedSource.asObservable();

    announceSession(loggedIn: boolean) {
        this.sessionAnnouncedSource.next(loggedIn);
    }

}

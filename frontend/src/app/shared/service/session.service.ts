import { Injectable } from '@angular/core';
import { Subject }    from 'rxjs/Subject';

@Injectable()
export class SessionService {

    private sessionAnnouncedSource = new Subject<boolean>();

    sessionAnnounced$ = this.sessionAnnouncedSource.asObservable();

    announceSession(loggedIn: boolean) {
        this.sessionAnnouncedSource.next(loggedIn);
    }

}

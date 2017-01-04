import {Component} from "@angular/core";
import {SessionService} from "./shared/service/session.service";

@Component({
    selector: 'p4a-root',
    templateUrl: './app.component.html',
    styleUrls: ['./app.component.css'],
    providers: [SessionService]
})

export class AppComponent {
    title = 'p4a works!';
}


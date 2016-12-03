import {Component, OnInit} from "@angular/core";
import {CtemplateService} from "ctemplate/service/ctemplate.service";
import {Ctemplate} from "ctemplate/ctemplate";

@Component({
    selector: 'ctemplate',
    templateUrl: 'html/ctemplate.component.html',
    styleUrls: ['ctemplate.component.css'],
    providers: [CtemplateService]
})

export class CtemplateComponent implements OnInit {
    ctemplates: Ctemplate[];
    parentCtemplate: Ctemplate;
    errorMsg: string;

    constructor(private ctemplateService: CtemplateService) {

    }

    ngOnInit(): void {
        this.initCtemplates();
    }

    initCtemplates(): void {
        this.ctemplateService.getCtemplates().subscribe(
                data => this.ctemplates = data,
                error => this.errorMsg = <any>error
        );
    }

    onDelete(ctemplate: Ctemplate): void {
        this.ctemplateService.removeCtemplate(ctemplate.id).subscribe();
        this.parentCtemplate = undefined;
    }

    onShow(ctemplate: Ctemplate): void {
        this.parentCtemplate = ctemplate;
    }

}

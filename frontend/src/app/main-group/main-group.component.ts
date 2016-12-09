import {Component, OnInit, Input, SimpleChanges} from "@angular/core";
import {MainGroupService} from "main-group/service/main-group.service";
import {MainGroup} from "main-group/main-group";
import {Ctemplate} from "../ctemplate/ctemplate";

@Component({
    selector: 'main-group',
    templateUrl: 'html/main-group.component.html',
    styleUrls: ['main-group.component.css'],
    providers: [MainGroupService]
})

export class MainGroupComponent implements OnInit {
    mainGroups: MainGroup[];
    parentMainGroup: MainGroup;
    objType: string;
    errorMsg: string;

    @Input() parentCtemplate: Ctemplate;

    constructor(private mainGroupService: MainGroupService) {
        this.objType = "MainGroup"
    }

    ngOnChanges(changes: SimpleChanges): void {
        for (let i in changes) {
            this.initMainGroups(changes[i].currentValue.id);
            this.parentMainGroup = undefined;
        }
    }

    ngOnInit(): void {
        this.initMainGroups(this.parentCtemplate.id);
    }

    initMainGroups(ctemplateId: number): void {
        this.mainGroupService.getSomeMainGroups("idTemplate", ctemplateId)
                .subscribe(data => this.mainGroups = data,
                error => this.errorMsg = <any>error
        );
    }

    onDelete(mainGroup: MainGroup): void {
        this.mainGroupService.removeMainGroup(mainGroup.id).subscribe();
        this.parentMainGroup = undefined;
    }

    onShow(mainGroup: MainGroup): void {
        this.parentMainGroup = mainGroup;
    }

}

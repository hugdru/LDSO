import {Component, OnInit, Input, SimpleChanges} from "@angular/core";
import {MainGroupService} from "main-group/service/main-group.service";
import {MainGroup} from "main-group/main-group";
import {AuditTemplate} from "../audit-template/audit-template";

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

    @Input() parentAuditTemplate: AuditTemplate;

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
        this.initMainGroups(this.parentAuditTemplate.id);
    }

    initMainGroups(auditTemplateId: number): void {
        this.mainGroupService.getSomeMainGroups("idTemplate", auditTemplateId)
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

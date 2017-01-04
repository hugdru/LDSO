import {Component, OnInit, Output, EventEmitter} from "@angular/core";
import {MainGroup} from "main-group/main-group";
import {SubGroup} from "../sub-group/sub-group";
import {MainGroupService} from "main-group/service/main-group.service";
import {SubGroupService} from "../sub-group/service/sub-group.service";
import {AuditTemplateService} from "../audit-template/service/audit-template.service";
import {AuditTemplate} from "../audit-template/audit-template";

@Component({
    selector: 'audit-select',
    templateUrl: './html/audit-select.component.html',
    styleUrls: ['./audit.component.css'],
    providers: [AuditTemplateService, MainGroupService, SubGroupService]
})

export class AuditSelectComponent implements OnInit {

    auditTemplate: AuditTemplate;
    mainGroups: MainGroup[];
    subGroups: SubGroup[];
    errorMsg: string;
    selectedSubGroups: SubGroup[];
    @Output() onDone = new EventEmitter<SubGroup[]>();

    constructor(private auditTemplateService: AuditTemplateService,
                private mainGroupService: MainGroupService,
                private subGroupService: SubGroupService) {
    }

    ngOnInit(): void {
        this.initAuditTemplate();
        this.mainGroups = [];
        this.subGroups = [];
        this.selectedSubGroups = [];
    }

    initAuditTemplate() : void {
        this.auditTemplateService.getCurrentAuditTemplate().subscribe(
                data => this.auditTemplate = data,
                error => this.errorMsg = <any> error
        );
    }

    showMainGroups(auditTemplate: AuditTemplate): void {
        this.initMainGroups(auditTemplate);
    }

    initMainGroups(auditTemplate: AuditTemplate): void {
        this.mainGroupService.getSomeMainGroups("idTemplate",
                auditTemplate.id).subscribe(
                data => this.mainGroups = data,
                error => this.errorMsg = <any> error
        );
    }

    showSubGroups(mainGroup: MainGroup): void {
        this.initSubGroups(mainGroup);
    }

    initSubGroups(mainGroup: MainGroup): void {
        this.subGroupService.getSomeSubGroups("idMaingroup",
                mainGroup.id).subscribe(data => this.subGroups = data);
    }

    toggleSubGroup(subGroup: SubGroup): void {
        var index = this.selectedSubGroups.map(
                function (x) {
                    return x.id;
                }).indexOf(subGroup.id);
        if (index > -1) {
            this.selectedSubGroups.splice(index, 1);
        }
        else {
            this.selectedSubGroups.push(subGroup);
        }
    }

    checkedSubGroup(subGroup: SubGroup): boolean {
        var index = this.selectedSubGroups.map(
                function (x) {
                    return x.id;
                }).indexOf(subGroup.id);
        return index > -1;
    }

    pressed(): void {
        this.onDone.emit(this.selectedSubGroups);
    }

}

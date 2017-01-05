import {Component, OnInit} from "@angular/core";
import {ActivatedRoute} from "@angular/router";
import {Property} from "property/property";
import {PropertyService} from "property/service/property.service";
import {SubGroup} from "../sub-group/sub-group";


@Component({
    selector: 'audit',
    templateUrl: './html/audit.component.html',
    styleUrls: ['./audit.component.css'],
    providers: [PropertyService]
})

export class AuditComponent implements OnInit {
    property: Property = new Property();
    selectedSubGroups: SubGroup[];
	auditId: number;

    constructor(private route: ActivatedRoute) {
    }

    ngOnInit(): void {
        this.property.id = +this.route.snapshot.params['id'];
    }

    onDone(selectedSubGroups: SubGroup[]): void {
        this.selectedSubGroups = selectedSubGroups;
    }

    sendId(auditId: number): void {
        this.auditId = auditId;
    }
}

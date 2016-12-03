import {
        Component,
        OnInit,
        OnChanges,
        SimpleChanges,
        Input
} from "@angular/core";
import {SubGroupService} from "sub-group/service/sub-group.service";
import {SubGroup} from "sub-group/sub-group";
import {MainGroup} from "main-group/main-group";

@Component({
    selector: 'sub-group',
    templateUrl: './html/sub-group.component.html',
    styleUrls: ['../main-group/main-group.component.css'],
    providers: [SubGroupService]
})

export class SubGroupComponent implements OnInit, OnChanges {
    subGroups: SubGroup[];
    parentSubGroup: SubGroup;

    @Input() parentMainGroup: MainGroup;

    constructor(private subGroupService: SubGroupService) {
    }

    ngOnChanges(changes: SimpleChanges): void {
        for (let i in changes) {
            this.initSubGroups(changes[i].currentValue.id);
            this.parentSubGroup = undefined;
        }
    }

    ngOnInit() {
        this.initSubGroups(this.parentMainGroup.id);
    }

    initSubGroups(mainGroupId: number): void {
        this.subGroupService
                .getSomeSubGroups("idMaingroup", mainGroupId)
                .subscribe(data => this.subGroups = data);

    }

    onDelete(subGroup: SubGroup): void {
        this.subGroupService.removeSubGroup(subGroup.id).subscribe();
        this.parentSubGroup = undefined;
    }

    onShow(subGroup: SubGroup): void {
        this.parentSubGroup = subGroup;
    }

}

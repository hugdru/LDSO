import {Component, Input, Output, EventEmitter, OnInit} from "@angular/core";
import {AccessibilityService} from "accessibility/service/accessibility.service";
import {Accessibility} from "accessibility/accessibility";

@Component({
    selector: 'accessibility-edit',
    templateUrl: '../ctemplate/html/ctemplate-edit.component.html',
    styleUrls: ['../ctemplate/ctemplate-edit.component.css'],
    providers: [AccessibilityService]
})

export class AccessibilityEditComponent implements OnInit {
    backupAccessibility: Accessibility;

    @Input() selectedObject: Accessibility;
    @Input() weight: number;
    @Output() onAction = new EventEmitter();

    constructor(private accessibilityService: AccessibilityService) {

    }

    ngOnInit(): void {
        this.backupAccessibility = new Accessibility();
        this.backupAccessibility.name = this.selectedObject.name;
        this.backupAccessibility.weight = this.selectedObject.weight;
    }

    pressed(updatedAccessibility: Accessibility): void {
        if (updatedAccessibility) {
            this.updateAccessibility();
        } else {
            this.selectedObject.name = this.backupAccessibility.name;
            this.selectedObject.weight = this.backupAccessibility.weight;
        }
        this.onAction.emit();
    }

    updateAccessibility(): void {
        this.accessibilityService.updateAccessibility(this.selectedObject)
                .subscribe();
    }

    checkPercentage(): boolean {
        return this.selectedObject.weight + this.weight > 100;
    }

}

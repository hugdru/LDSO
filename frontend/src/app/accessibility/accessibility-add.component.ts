import {Component, Input, Output, EventEmitter, OnInit} from "@angular/core";
import {AccessibilityService} from "accessibility/service/accessibility.service";
import {Accessibility} from "accessibility/accessibility";
import {Criterion} from "criterion/criterion";

@Component({
    selector: 'accessibility-add',
    templateUrl: '../ctemplate/html/ctemplate-edit.component.html',
    styleUrls: ['../ctemplate/ctemplate-edit.component.css'],
    providers: [AccessibilityService]
})

export class AccessibilityAddComponent implements OnInit {
    selectedObject: Accessibility;

    @Input() objType: string;
    @Input() criterion: Criterion;
    @Input() weight: number;
    @Output() onAdd = new EventEmitter<Accessibility>();

    constructor(private accessibilityService: AccessibilityService) {

    }

    ngOnInit(): void {
        this.selectedObject = new Accessibility();
    }

    pressed(newAccessibility: Accessibility): void {
        if (newAccessibility) {
            this.addAccessibility();
        }
        this.onAdd.emit(newAccessibility);
    }

    addAccessibility(): void {
        this.accessibilityService
                .setAccessibility(this.selectedObject , this.criterion.id)
                .subscribe(
                        response => this.selectedObject.id = response.json().id
                );
    }

    checkPercentage(): boolean {
        return this.selectedObject.weight + this.weight != 100;
    }
}

import {Component, Input, Output, EventEmitter, OnInit} from "@angular/core";
import {AccessibilityService} from "accessibility/service/accessibility.service";
import {Accessibility} from "accessibility/accessibility";
import {Criterion} from "criterion/criterion";

@Component({
    selector: 'accessibility-add',
    templateUrl: '../main-group/html/main-group-edit.component.html',
    styleUrls: ['../main-group/main-group-edit.component.css'],
    providers: [AccessibilityService]
})

export class AccessibilityAddComponent implements OnInit {
    selectedObject: Accessibility;

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
        this.selectedObject.criterion = this.criterion._id;
        this.accessibilityService.setAccessibility(this.selectedObject)
                .subscribe(
                        response => this.selectedObject._id = response.json()
                );
    }

    checkPercentage(): boolean {
        return this.selectedObject.weight + this.weight > 100;
    }
}

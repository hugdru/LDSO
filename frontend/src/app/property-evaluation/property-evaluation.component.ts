import { Component, OnInit ,Input} from '@angular/core';
import {PropertyEvaluation} from "./property-evaluation";
import {PropertyEvaluationService} from "./services/porpertyEvaluation.service";

@Component({
  selector: 'app-property-evaluation',
  templateUrl: 'html/property-evaluation.component.html',
  styleUrls: ['./property-evaluation.component.css']
})
export class PropertyEvaluationComponent implements OnInit {
    propertyE: PropertyEvaluation;
    @Input() propertyid: number;

    constructor(private propertyEvaluationS: PropertyEvaluationService) {
    }

    ngOnInit(): void {
        this.propertyEvaluationS.getPropertyEvaluation(this.propertyid)
            .subscribe(data => this.propertyE = data);

    }
}

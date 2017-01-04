import { Component, OnInit ,Input} from '@angular/core';
import {PropertyEvaluation} from "./property-evaluation";
import {PropertyEvaluationService} from "./service/porpertyEvaluation.service";

@Component({
  selector: 'app-property-evaluation',
  templateUrl: 'html/property-evaluation.component.html',
  styleUrls: ['./property-evaluation.component.css'],
  providers: [PropertyEvaluationService]
})
export class PropertyEvaluationComponent implements OnInit {
    propertyE: PropertyEvaluation;
    @Input() propertyid: number;
    errorMsg: string;

    constructor(private PorpertyEv: PropertyEvaluationService) {

        this.PorpertyEv.getPropertyEvaluation(1).subscribe(

            response=> console.log(response)
            ,
            error => this.errorMsg = <any>error

        );

       // console.log(this.propertyE.rating);
    }

    ngOnInit(): void {

        //console.log(this.propertyE.rating);
    }

}

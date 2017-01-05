import { Component, OnInit ,Input} from '@angular/core';
import {PropertyEvaluation} from "./property-evaluation";
import {PropertyEvaluationService} from "./service/porpertyEvaluation.service";
import {ActivatedRoute} from '@angular/router';
import {Subscription} from "rxjs";

@Component({
  selector: 'app-property-evaluation',
  templateUrl: 'html/property-evaluation.component.html',
  styleUrls: ['./property-evaluation.component.css'],
  providers: [PropertyEvaluationService]
})
export class PropertyEvaluationComponent implements OnInit {
    propertyE:  PropertyEvaluation[];
    errorMsg: string;
    private subcription: Subscription;
    id: string;


   constructor(private PorpertyEv: PropertyEvaluationService, private activatedRoute: ActivatedRoute) {
       this.subcription =activatedRoute.params.subscribe(
           (param: any) => this.id = param['id']
       );
        this.PorpertyEv.getPropertyEvaluation(1).subscribe(

            data=> this.propertyE = data

            ,
            error => this.errorMsg = <any>error

        );
        console.log(this.propertyE);
       console.log(+this.id);

    }

    propertyMock =
    [
        {id: 1, name: "casa das francesinhas",adress:"rua x do catino y da freguesia z",imagePath:"http://cdn1.buuteeq.com/upload/2020657/foto3.jpg.694x0_default.jpg"}
    ]
    ;
    auditorMock =
    [
        {id: 1, name: "carlos chiquinho"},
        {id: 2, name: "edudarto gomes"},
        {id: 3, name: "Miguel costa"},
        {id: 4, name: "osvaldo garcia"}
    ]
    ;
    propertyEvalutationEMock =
    [
        {id: 1, property:this.propertyMock[0],  auditor: this.auditorMock[0] ,idTemplate:1,rating:31,createdDate:"0",finishedDate:"3",coment:"Scams in this decade are very prevalent, unfortunately preying on the kind and naive. The blog above is unfortunately pretty accurate. The CNN article I found on a more general note had an entire list of the worst 50 charities distinguished by the percentage falsely promised to go to the cause. This link [http://www.tampabay.com/americas-worst-charities/ ] will take you directly to this list."},
        {id: 2, property:this.propertyMock[0],  auditor: this.auditorMock[1] ,idTemplate:1,rating:33,createdDate:"0",finishedDate:"4",coment:"Using state and federal records, the Times and CIR identified nearly 6,000 charities that have chosen to pay for-profit companies to raise their donations……The 50 worst charities in America devote less than 4% of donations raised to direct cash aid. Some charities gave even less. Over a decade, one diabetes charity raised nearly $14 million and gave about $10,000 to patients. Six spent no cash at all on their cause"},
        {id: 3, property:this.propertyMock[0],  auditor: this.auditorMock[2] ,idTemplate:1,rating:25,createdDate:"0",finishedDate:"3",coment:"Scams in this decade are very prevalent, unfortunately preying on the kind and naive. The blog above is unfortunately pretty accurate. The CNN article I found on a more general note had an entire list of the worst 50 charities distinguished by the percentage falsely promised to go to the cause. This link [http://www.tampabay.com/americas-worst-charities/ ] will take you directly to this list."},
    ]
    ;



    ngOnInit(): void {


    }

}

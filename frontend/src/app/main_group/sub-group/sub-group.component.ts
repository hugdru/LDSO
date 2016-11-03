import {Component, OnDestroy} from '@angular/core';
import {Router, ActivatedRoute} from '@angular/router';
import {FormGroup, FormControl, Validators, FormArray} from "@angular/forms";
import {Subscription} from "rxjs";
import {FormInfoService} from './../service/form-info.service';

@Component({
  selector: 'app-sub-group',
  templateUrl: './sub-group.component.html',
  styleUrls: ['./sub-group.component.css'],
  providers: [FormInfoService]
})
export class SubGroupComponent implements OnDestroy{
  private subcription: Subscription;

  subFormGroup: FormGroup;


  id: string;
  constructor(private router: Router, private activatedRoute: ActivatedRoute, private formVar: FormInfoService) {
    this.subcription =activatedRoute.params.subscribe(
      (param: any) => this.id = param['id']
    );

    this.subFormGroup = formVar.getFormGroup();
  }
  ngOnDestroy(){
    this.subcription.unsubscribe();
  }
  onSubmit(){
    console.log(this.id);
  }
}

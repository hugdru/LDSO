import {Component, OnDestroy} from '@angular/core';
import {Router, ActivatedRoute} from '@angular/router'
import {Subscription} from "rxjs";
import {FormGroup, FormControl, Validators, FormArray} from "@angular/forms";
import {FormInfoService} from './service/form-info.service';

@Component({
  selector: 'app-main_group',
  templateUrl: 'mainGroup.component.html',
  styleUrls: ['mainGroup.component.css'],
  providers: [FormInfoService]
})

export class MainGroupComponent implements OnDestroy {
  private subcription: Subscription;
  myForm : FormGroup;

  id: string;
  constructor(private router: Router, private activatedRoute: ActivatedRoute, private formVar: FormInfoService) {
    this.subcription =activatedRoute.params.subscribe();
    this.myForm = formVar.getFormGroup();
  }

  onSubmit(){
    console.log(this.myForm.value);

  }

  ngOnDestroy(){
    this.subcription.unsubscribe();
  }
}

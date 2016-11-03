import {Component, OnDestroy} from '@angular/core';
import {Router, ActivatedRoute} from '@angular/router'
import {Subscription} from "rxjs";
import {FormGroup, FormControl, Validators, FormArray} from "@angular/forms";
import {FormInfoService} from './service/form-info.service';

@Component({
  selector: 'app-main_group',
  templateUrl: 'main_group.component.html',
  styleUrls: ['main_group.component.css'],
  providers: [FormInfoService]
})

export class Main_groupComponent implements OnDestroy {
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

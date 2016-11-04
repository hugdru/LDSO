import {Component, OnDestroy} from '@angular/core';
import {Router, ActivatedRoute} from '@angular/router';
import {FormGroup, FormControl, Validators, FormArray} from "@angular/forms";
import {Subscription} from "rxjs";
import {FormInfoService} from '../mainGroup/service/form-info.service';

@Component({
  selector: 'app-sub-group',
  templateUrl: './html/sub-group.component.html',
  styleUrls: ['subGroup.component.css'],
  providers: [FormInfoService]
})
export class SubGroupComponent implements OnDestroy{
  private subcription: Subscription;
  formGroup: FormGroup;

  id: string;
  constructor(private router: Router, private activatedRoute: ActivatedRoute, private formVar: FormInfoService) {
    this.subcription =activatedRoute.params.subscribe(
      (param: any) => this.id = param['id']
    );
    this.formGroup = formVar.getFormGroup();
    //this.subFormGroup = formVar.getSubFormGroup(this.id);
  }
  ngOnDestroy(){
    this.subcription.unsubscribe();
  }
  onSubmit(){
   // console.log(this.id);
   // console.log(this.formGroup.value);
   // console.log(this.subFormGroup.value);
  }
  onAddSubgroup(SubGroup_Name:string){

    this.formVar.addSubGroup(this.id,SubGroup_Name);
   // this.formVar.addCriterios(this.id,SubGroup_Name,'teste');
    console.log(SubGroup_Name);
    console.log(this.formGroup.value[this.id]);
  }
}

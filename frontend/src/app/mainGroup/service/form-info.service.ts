import { Injectable } from '@angular/core';
import {FormGroup, FormControl, Validators, FormArray} from "@angular/forms";

@Injectable()
export class FormInfoService{
  myFormSubForm: FormGroup;

  constructor() {
    this.myFormSubForm = new FormGroup(
      {
          'acessos' : new FormArray([]),
          'bens_e_servicos' : new FormArray([]),
          'percurso_exterior' : new FormArray([]),
          'percurso_interior' : new FormArray([])
      }
    );
  }

  addSubGroup(main_Group: string, sub_group: string){
    this.myFormSubForm.value[main_Group].push({[sub_group]: new FormGroup({})});
  }
/*
  addCriterios(main_Group: string, sub_group: string, criteria: string){
    (this.myFormSubForm.value[main_Group]).value[sub_group].push(new FormGroup(
      {criteria : new FormControl(['bom',Validators.required()]),
        'fisica': new FormControl(['0',Validators.required()]),
        'auditiva': new FormControl(['0',Validators.required()]),
        'visual': new FormControl(['0',Validators.required()]),
        'cognitiva': new FormControl(['0',Validators.required()])
      }
    ));
  }*/
  getFormGroup(){
    return this.myFormSubForm;
  }

  getSubFormGroup(main_Group: string){
    return this.myFormSubForm.value[main_Group];
  }
}

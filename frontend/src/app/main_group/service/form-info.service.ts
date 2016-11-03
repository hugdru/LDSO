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
    this.myFormSubForm.value[main_Group].push(new FormGroup(
                                                {sub_group: new FormArray([])}
                                              ));
  }

  addCriterios(main_Group: string, sub_group: string, criteria: string){
    (this.myFormSubForm.value[main_Group]).value[sub_group].push(new FormGroup(
      {criteria : new FormControl([])}
    ));
  }



  getFormGroup(){
    return this.myFormSubForm;
  }
}

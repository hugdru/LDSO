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
    console.log(this.myFormSubForm.value[main_Group]);
    this.myFormSubForm.value[main_Group].push(
                        {'subGroup' : [sub_group],
                         'criteria': new FormArray([])
                        });
  }
  addCriterios(main_Group: string, sub_group: string, new_criteria: string){
    for(let group of this.myFormSubForm.value[main_Group]){
      if(group.subGroup === sub_group){
         group.criteria.controls.push(
            {
              'criteria_name' : [new_criteria],
              'weigths' :
                {
                  'Physical' : '20',
                  'Auditor' : '20',
                  'Visual' : '20',
                  'Cognitive' : '20',
                  'GoodPractice' : 'no'
                },
              'Comments': 'write coments'
            }
          );
      }
    }
  }

  getFormGroup(){
    return this.myFormSubForm;
  }

  getSubFormGroup(main_Group: string){
    return this.myFormSubForm.value[main_Group];
  }
}

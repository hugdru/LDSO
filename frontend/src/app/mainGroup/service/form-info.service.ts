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
  addCriterios(main_Group: string, sub_group: string, criteria: string){
    this.myFormSubForm.value[main_Group][sub_group].push(new FormGroup(
      {
        [criteria]:new FormControl(['bom']),
        'physical':new FormControl(['0']),
        'auditory':new FormControl(['0']),
        'visual':new FormControl(['0']),
        'cognitive':new FormControl(['0']),
        'comentarios':new FormControl(['Escrever comentarios'])
      }
    ));
  }

  getFormGroup(){
    return this.myFormSubForm;
  }

  getSubFormGroup(main_Group: string){
    return this.myFormSubForm.value[main_Group];
  }
}

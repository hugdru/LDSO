import { Component, OnInit } from '@angular/core';
import {FormGroup, FormControl, Validators, FormArray} from "@angular/forms";

@Component({
  selector: 'app-criterion',
  templateUrl: './criterion.component.html',
  styleUrls: ['./criterion.component.css']
})
export class CriterionComponent implements OnInit {
  myForm: FormGroup;
  constructor() {

  }

  ngOnInit() {
  }

  onSubmit(){
    //console.log(this.myForm.value);

  }
  }

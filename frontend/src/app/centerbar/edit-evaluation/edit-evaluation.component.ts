import { Component, OnInit } from '@angular/core';

@Component({
  selector: 'p4a-edit-evaluation',
  templateUrl: './edit-evaluation.component.html',
  styleUrls: ['./edit-evaluation.component.css']
})

export class EditEvaluationComponent implements OnInit {

  private items = [ ["acesso",["elevador","rampa"]],
                    ["percurso exterior",["passeios","autocarros"]]
                  ];

  private item = ["teste",["teste1","teste2"]];

  addItem() {
	  this.items.push(this.item);
  }

  constructor() { }

  ngOnInit() {
  }

}

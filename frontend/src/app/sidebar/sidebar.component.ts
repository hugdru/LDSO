import { Component, OnInit , Input}  from '@angular/core';

@Component({
  selector: 'p4a-sidebar',
  templateUrl: './sidebar.component.html',
  styleUrls: ['./sidebar.component.css']
})

export class SidebarComponent implements OnInit {

  @Input() Recipe = {
    name:"Auditoria",
    number:27
  };

  constructor() { }

  ngOnInit() {
  }
}

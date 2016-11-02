import {Component, OnDestroy} from '@angular/core';
import {Router, ActivatedRoute} from '@angular/router'
import {Subscription} from "rxjs";

@Component({
  selector: 'app-sub-group',
  templateUrl: './sub-group.component.html',
  styleUrls: ['./sub-group.component.css']
})
export class SubGroupComponent implements OnDestroy{
  private subcription: Subscription;

  id: string;
  constructor(private router: Router, private activatedRoute: ActivatedRoute) {
    this.subcription =activatedRoute.params.subscribe(
      (param: any) => this.id = param['id']


    );
  }
  ngOnDestroy(){
    this.subcription.unsubscribe();
  }
  onPress(){
    console.log(this.id);
  }
}

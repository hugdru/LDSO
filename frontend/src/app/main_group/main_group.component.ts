import {Component, OnDestroy} from '@angular/core';
import {Router, ActivatedRoute} from '@angular/router'
import {Subscription} from "rxjs";


@Component({
  selector: 'app-main_group',
  templateUrl: 'main_group.component.html',
  styleUrls: ['main_group.component.css']
})
export class Main_groupComponent implements OnDestroy {
  private subcription: Subscription;

  id: string;
  constructor(private router: Router, private activatedRoute: ActivatedRoute) {
    this.subcription =activatedRoute.params.subscribe(

    );
  }
  ngOnDestroy(){
    this.subcription.unsubscribe();
  }
}

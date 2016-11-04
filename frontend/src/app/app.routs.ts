import {Routes, RouterModule} from '@angular/router';

import {MainGroupComponent} from './mainGroup/mainGroup.component';
import {SubGroupComponent} from './subGroup/subGroup.component';

const APP_ROUTES: Routes =[
  {path: 'mainGroup',component: MainGroupComponent},
  {path: 'subGroup/:id',component: SubGroupComponent}
];

export const routing = RouterModule.forRoot(APP_ROUTES);

import {Routes, RouterModule} from '@angular/router';

import {MainGroupComponent} from './main-group/main-group.component';
import {SubGroupComponent} from './sub-group/sub-group.component';

const APP_ROUTES: Routes =[
  {path: 'main-group',component: MainGroupComponent},
  {path: 'sub-group/:id',component: SubGroupComponent}
];

export const routing = RouterModule.forRoot(APP_ROUTES);

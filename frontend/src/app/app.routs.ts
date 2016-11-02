import {Routes, RouterModule} from '@angular/router';

import {Main_groupComponent} from './main_group/main_group.component';
import {SubGroupComponent} from './main_group/sub-group/sub-group.component';

const APP_ROUTES: Routes =[
  {path: 'main_group',component: Main_groupComponent},
  {path: 'sub_group/:id',component: SubGroupComponent}
];

export const routing = RouterModule.forRoot(APP_ROUTES);

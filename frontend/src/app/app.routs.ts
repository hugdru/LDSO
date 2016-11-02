import {Routes, RouterModule} from '@angular/router';

import {CriterionComponent} from './criterion/criterion.component';

const APP_ROUTES: Routes =[
  {path: 'criterion',component: CriterionComponent}
];

export const routing = RouterModule.forRoot(APP_ROUTES);

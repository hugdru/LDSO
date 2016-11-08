import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { MainGroupComponent } from 'main-group/main-group.component';
import { SubGroupComponent } from 'sub-group/sub-group.component';

const routes: Routes = [
	{
		path: '',
		redirectTo: '/main-group', //change to homepage once implemented
		pathMatch: 'full'
	},
	{
		path: 'main-group',
		component: MainGroupComponent
	}
];

@NgModule({
	imports: [ RouterModule.forRoot(routes) ],
	exports: [ RouterModule ]
})

export class AppRoutingModule {

}


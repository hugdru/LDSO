import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { MainGroupComponent } from 'main-group/main-group.component';
import { AuditComponent } from 'audit/audit.component';

const routes: Routes = [
	{
		path: '',
		redirectTo: '/main-group', //change to homepage once implemented
		pathMatch: 'full'
	},
	{
		path: 'main-group',
		component: MainGroupComponent
	},
	{
		path: 'audit/:id',
		component: AuditComponent
	}
];

@NgModule({
	imports: [ RouterModule.forRoot(routes) ],
	exports: [ RouterModule ]
})

export class AppRoutingModule {

}


import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { MainGroupComponent } from 'main-group/main-group.component';
import { AuditComponent } from 'audit/audit.component';
import {ListPropertiesComponent} from "./property/list-properties.component";

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
	},
	{
		path: 'listProperties',
		component: ListPropertiesComponent
	}
];

@NgModule({
	imports: [ RouterModule.forRoot(routes) ],
	exports: [ RouterModule ]
})

export class AppRoutingModule {

}


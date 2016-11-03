import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import {FormsModule, ReactiveFormsModule} from '@angular/forms';
import { HttpModule } from '@angular/http';

import { AppComponent } from './app.component';
import { HeaderComponent } from './header/header.component';
import { SidebarComponent } from './sidebar/sidebar.component';
import { FootbarComponent } from './footbar/footbar.component';
import { CenterComponent } from './center/center.component';
import { AuditComponent } from './audit/audit.component';
import { EditEvaluationComponent } from './audit/edit-evaluation.component';
import { Main_groupComponent } from './main_group/main_group.component';

import { routing } from './app.routs';
import { SubGroupComponent } from './main_group/sub-group/sub-group.component';


@NgModule({
	declarations: [
		AppComponent,
		HeaderComponent,
		SidebarComponent,
		FootbarComponent,
		CenterComponent,
		AuditComponent,
		EditEvaluationComponent,
    Main_groupComponent,
    SubGroupComponent
	],
	imports: [
		BrowserModule,
		FormsModule,
		HttpModule,
    routing,
    ReactiveFormsModule
	],
	providers: [],
	bootstrap: [AppComponent]
})

export class AppModule {

}

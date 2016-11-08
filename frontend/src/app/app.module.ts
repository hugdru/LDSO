import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { HttpModule } from '@angular/http';

import { AppComponent } from 'app.component';
import { HeaderComponent } from 'header/header.component';
import { SidebarComponent } from 'sidebar/sidebar.component';
import { FootbarComponent } from 'footbar/footbar.component';
import { CenterComponent } from 'center/center.component';
import { AuditComponent } from 'audit/audit.component';
import { EditEvaluationComponent } from 'audit/edit-evaluation.component';
import { MainGroupComponent } from 'main-group/main-group.component';
import { SubGroupComponent } from 'sub-group/sub-group.component';

import { HandlerService } from 'handler.service';

import { routing } from 'app.routs';


@NgModule({
	declarations: [
		AppComponent,
		HeaderComponent,
		SidebarComponent,
		FootbarComponent,
		CenterComponent,
		AuditComponent,
		EditEvaluationComponent,
		MainGroupComponent,
		SubGroupComponent
	],
	imports: [
		BrowserModule,
		FormsModule,
		HttpModule,
		routing,
		ReactiveFormsModule
	],
	providers: [ HandlerService ],
	bootstrap: [ AppComponent ]
})

export class AppModule {

}

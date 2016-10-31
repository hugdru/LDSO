import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { HttpModule } from '@angular/http';

import { AppComponent } from './app.component';
import { HeaderComponent } from './header/header.component';
import { SidebarComponent } from './sidebar/sidebar.component';
import { FootbarComponent } from './footbar/footbar.component';
import { CenterComponent } from './center/center.component';
import { AuditComponent } from './audit/audit.component';
import { EditEvaluationComponent } from './audit/edit-evaluation.component';
import { SharedData } from './shared/shared-data';
import { Group } from './group/group';
import { SubGroup } from './sub-group/sub-group';

@NgModule({
	declarations: [
		AppComponent,
		HeaderComponent,
		SidebarComponent,
		FootbarComponent,
		CenterComponent,
		AuditComponent,
		EditEvaluationComponent,
		SharedData,
		Group
	],
	imports: [
		BrowserModule,
		FormsModule,
		HttpModule
	],
	providers: [],
	bootstrap: [AppComponent]
})

export class AppModule {

}

import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { HttpModule } from '@angular/http';
import { FormsModule } from '@angular/forms';

import { NgbModule } from '@ng-bootstrap/ng-bootstrap';
import { AppRoutingModule } from 'app-routing.module';

import { AppComponent } from 'app.component';
import { HeaderComponent } from 'header/header.component';
import { SidebarComponent } from 'sidebar/sidebar.component';
import { FootbarComponent } from 'footbar/footbar.component';
import { CenterComponent } from 'center/center.component';
import { MainGroupComponent } from 'main-group/main-group.component';
import { MainGroupEditComponent } from 'main-group/main-group-edit.component';
import { MainGroupAddComponent } from 'main-group/main-group-add.component';
import { SubGroupComponent } from 'sub-group/sub-group.component';
import { SubGroupEditComponent } from 'sub-group/sub-group-edit.component';
import { SubGroupAddComponent } from 'sub-group/sub-group-add.component';
import { ListManageComponent } from 'shared/list-manage/list-manage.component';

import { HandlerService } from 'handler.service';

@NgModule({
	declarations: [
		AppComponent,
		HeaderComponent,
		SidebarComponent,
		FootbarComponent,
		CenterComponent,
		MainGroupComponent,
		MainGroupEditComponent,
		MainGroupAddComponent,
		SubGroupComponent,
		SubGroupEditComponent,
		SubGroupAddComponent,
		ListManageComponent
	],
	imports: [
		AppRoutingModule,
		BrowserModule,
		HttpModule,
		FormsModule,
		NgbModule.forRoot()
	],
	providers: [ HandlerService ],
	bootstrap: [ AppComponent ]
})

export class AppModule {

}

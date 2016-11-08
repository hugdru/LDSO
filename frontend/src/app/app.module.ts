import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { HttpModule } from '@angular/http';
import { FormsModule } from '@angular/forms';

import { AppRoutingModule } from 'app-routing.module';

import { AppComponent } from 'app.component';
import { HeaderComponent } from 'header/header.component';
import { SidebarComponent } from 'sidebar/sidebar.component';
import { FootbarComponent } from 'footbar/footbar.component';
import { CenterComponent } from 'center/center.component';
import { MainGroupComponent } from 'main-group/main-group.component';
import { SubGroupComponent } from 'sub-group/sub-group.component';

import { HandlerService } from 'handler.service';

@NgModule({
	declarations: [
		AppComponent,
		HeaderComponent,
		SidebarComponent,
		FootbarComponent,
		CenterComponent,
		MainGroupComponent,
		SubGroupComponent
	],
	imports: [
		AppRoutingModule,
		BrowserModule,
		HttpModule,
		FormsModule
	],
	providers: [ HandlerService ],
	bootstrap: [ AppComponent ]
})

export class AppModule {

}

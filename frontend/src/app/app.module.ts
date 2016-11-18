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
import { CriterionComponent } from 'criterion/criterion.component';
import { CriterionEditComponent } from 'criterion/criterion-edit.component';
import { CriterionAddComponent } from 'criterion/criterion-add.component';
import { AccessibilityComponent } from 'accessibility/accessibility.component';
import {
	AccessibilityEditComponent
} from 'accessibility/accessibility-edit.component';
import {
	AccessibilityAddComponent
} from 'accessibility/accessibility-add.component';
import { ListManageComponent } from 'shared/list-manage/list-manage.component';
import {
	ListSelectAddComponent
} from 'shared/list-select/list-select-add.component';
import {
	ListSelectEditComponent
} from 'shared/list-select/list-select-edit.component';
import {
	ListShowValueComponent
} from 'shared/list-show-value/list-show-value.component';
import {
	ListPropertiesComponent
} from 'property/list-properties.component';
import{
		PropertiesInfoComponent,
} from 'property/property-info.component';
import { AuditComponent } from 'audit/audit.component';
import { AuditListGroupsComponent } from 'audit/audit-list-groups.component';
import { AuditSelectComponent } from 'audit/audit-select.component';

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
		CriterionComponent,
		CriterionEditComponent,
		CriterionAddComponent,
		AccessibilityComponent,
		AccessibilityEditComponent,
		AccessibilityAddComponent,
		ListManageComponent,
		ListSelectAddComponent,
		ListSelectEditComponent,
		ListShowValueComponent,
		ListPropertiesComponent,
		PropertiesInfoComponent,
		AuditComponent,
		AuditListGroupsComponent,
		AuditSelectComponent
	],
	imports: [
		AppRoutingModule,
		BrowserModule,
		HttpModule,
		FormsModule,
		NgbModule.forRoot()
	],
	providers: [ HandlerService ],
	bootstrap: [ AppComponent]
})

export class AppModule {

}

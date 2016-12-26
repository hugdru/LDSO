import {BrowserModule} from "@angular/platform-browser";
import {NgModule} from "@angular/core";
import {HttpModule} from "@angular/http";
import {FormsModule} from "@angular/forms";
import {NgbModule} from "@ng-bootstrap/ng-bootstrap";
import {ModalModule} from "ng2-modal";
import {MaterialModule} from '@angular/material';
import {FileSelectDirective, FileDropDirective}
		from 'ng2-file-upload/ng2-file-upload';
import 'hammerjs';

import {AppRoutingModule} from "app-routing.module";
import {AppComponent} from "app.component";
import {HeaderComponent} from "header/header.component";
import {SidebarComponent} from "sidebar/sidebar.component";
import {FootbarComponent} from "footbar/footbar.component";
import {CenterComponent} from "center/center.component";
import {EqualValidator} from "shared/equal-validator.directive";
import {LoginComponent} from "login/login.component";
import {RegisterComponent} from "./register/register.component";
import {AuditTemplateComponent} from "audit-template/audit-template.component";
import {AuditTemplateEditComponent} from "audit-template/audit-template-edit.component";
import {AuditTemplateAddComponent} from "audit-template/audit-template-add.component";
import {MainGroupComponent} from "main-group/main-group.component";
import {MainGroupEditComponent} from "main-group/main-group-edit.component";
import {MainGroupAddComponent} from "main-group/main-group-add.component";
import {SubGroupComponent} from "sub-group/sub-group.component";
import {SubGroupEditComponent} from "sub-group/sub-group-edit.component";
import {SubGroupAddComponent} from "sub-group/sub-group-add.component";
import {CriterionComponent} from "criterion/criterion.component";
import {CriterionEditComponent} from "criterion/criterion-edit.component";
import {CriterionAddComponent} from "criterion/criterion-add.component";
import {AccessibilityComponent} from "accessibility/accessibility.component";
import {AccessibilityEditComponent} from "accessibility/accessibility-edit.component";
import {AccessibilityAddComponent} from "accessibility/accessibility-add.component";
import {ListManageComponent} from "shared/list-manage/list-manage.component";
import {ListSelectAddComponent} from "shared/list-select/list-select-add.component";
import {ListSelectEditComponent} from "shared/list-select/list-select-edit.component";
import {ListShowValueComponent} from "shared/list-show-value/list-show-value.component";
import {ShowListComponent} from "shared/show-list/show-list.component";
import {ListPropertiesComponent} from "property/list-properties.component";
import {PropertiesInfoComponent} from "property/property-info.component";
import {AuditComponent} from "audit/audit.component";
import {AuditEvaluateComponent} from "audit/audit-evaluate.component";
import {AuditEvaluateCriterionComponent} from "audit/audit-evaluate-criterion.component";
import {AuditSelectComponent} from "audit/audit-select.component";
import {RemarkAddComponent} from "remark/remark-add.component";
import {ImageSingleUploadComponent} from "images/image-single-upload.component";
import {HandlerService} from "handler.service";
import { MainPageComponent } from './main-page/main-page.component';

@NgModule({
    declarations: [
        FileSelectDirective,
        FileDropDirective,
        AppComponent,
        HeaderComponent,
        SidebarComponent,
        FootbarComponent,
        CenterComponent,
        LoginComponent,
        EqualValidator,
        RegisterComponent,
        AuditTemplateComponent,
        AuditTemplateEditComponent,
        AuditTemplateAddComponent,
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
        ShowListComponent,
        PropertiesInfoComponent,
        AuditComponent,
        AuditEvaluateComponent,
        AuditEvaluateCriterionComponent,
        AuditSelectComponent,
    	RemarkAddComponent,
		ImageSingleUploadComponent,
		MainPageComponent
    ],
    imports: [
        AppRoutingModule,
        BrowserModule,
        HttpModule,
        FormsModule,
        NgbModule.forRoot(),
        ModalModule,
		MaterialModule.forRoot()
    ],
    providers: [HandlerService],
    bootstrap: [AppComponent]
})

export class AppModule {

}

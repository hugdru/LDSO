import {NgModule} from "@angular/core";
import {RouterModule, Routes} from "@angular/router";
import {AuditTemplateComponent} from "audit-template/audit-template.component";
import {MainPageComponent} from "main-page/main-page.component";
import {AuditComponent} from "audit/audit.component";
import {ListPropertiesComponent} from "./property/list-properties.component";
import {LoginComponent} from "./login/login.component";
import {RegisterComponent} from "./register/register.component";
import { PropertyEvaluationComponent } from './property-evaluation/property-evaluation.component';

const routes: Routes = [
    {
        path: '',
        redirectTo: '/main-page', //change to homepage once implemented
        pathMatch: 'full'
    },
    {
        path: 'main-page',
        component: MainPageComponent
    },
    {
        path: 'auditTemplate',
        component: AuditTemplateComponent
    },
    {
        path: 'audit/:id',
        component: AuditComponent
    },
    {
        path: 'listProperties',
        component: ListPropertiesComponent
    },
    {
        path: 'login',
        component: LoginComponent },
    {
        path: 'register',
        component: RegisterComponent },
    {
        path: 'propertyEvaluation',
        component: PropertyEvaluationComponent
    },

];

@NgModule({
    imports: [RouterModule.forRoot(routes)],
    exports: [RouterModule]
})

export class AppRoutingModule {

}


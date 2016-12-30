import {NgModule} from "@angular/core";
import {RouterModule, Routes} from "@angular/router";
import {AuditTemplateComponent} from "audit-template/audit-template.component";
import {AuditComponent} from "audit/audit.component";
import {ListPropertiesComponent} from "./property/list-properties.component";
import {LoginComponent} from "./login/login.component";
import {RegisterComponent} from "./register/register.component";

const routes: Routes = [
    {
        path: '',
        redirectTo: '/auditTemplate', //change to homepage once implemented
        pathMatch: 'full'
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

];

@NgModule({
    imports: [RouterModule.forRoot(routes)],
    exports: [RouterModule]
})

export class AppRoutingModule {

}


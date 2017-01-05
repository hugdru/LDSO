import {NgModule, Injectable} from "@angular/core";
import {RouterModule, Routes} from "@angular/router";
import {AuditTemplateComponent} from "audit-template/audit-template.component";
import {MainPageComponent} from "main-page/main-page.component";
import {AuditComponent} from "audit/audit.component";
import {ListPropertiesComponent} from "./property/list-properties.component";
import {LoginComponent} from "./login/login.component";
import {RegisterComponent} from "./register/register.component";
import {PropertyEvaluationComponent} from "./property-evaluation/property-evaluation.component";
import {SuperAdminGuard} from "shared/service/superadmin.guard";
import {AuditGuard} from "./shared/service/audit.guard";

const routes: Routes = [
    {
        path: '',
        redirectTo: '/main-page',
        pathMatch: 'full'
    },
    {
        path: 'main-page',
        component: MainPageComponent
    },
    {
        path: 'auditTemplate',
        component: AuditTemplateComponent,
        canActivate: [SuperAdminGuard]
    },
    {
        path: 'audit/:id',
        component: AuditComponent,
        canActivate: [AuditGuard]
    },
    {
        path: 'listPropertiesAudit',
        component: ListPropertiesComponent,
        canActivate: [AuditGuard]
    },
    {
        path: 'listPropertiesView',
        component: ListPropertiesComponent
    },
    {
        path: 'login',
        component: LoginComponent
    },
    {
        path: 'register',
        component: RegisterComponent
    },
    {
        path: 'propertyEvaluation/:id',
        component: PropertyEvaluationComponent
    },

];

@NgModule({
    imports: [RouterModule.forRoot(routes)],
    exports: [RouterModule]
})

@Injectable()
export class AppRoutingModule {

}


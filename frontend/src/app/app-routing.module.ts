import {NgModule} from "@angular/core";
import {RouterModule, Routes} from "@angular/router";
import {CtemplateComponent} from "ctemplate/ctemplate.component";
import {AuditComponent} from "audit/audit.component";
import {ListPropertiesComponent} from "./property/list-properties.component";
import {LoginComponent} from "./login/login.component";

const routes: Routes = [
    {
        path: '',
        redirectTo: '/ctemplate', //change to homepage once implemented
        pathMatch: 'full'
    },
    {
        path: 'ctemplate',
        component: CtemplateComponent
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
/*    {
        path: 'signup',
        component: SignupComponent },*/

];

@NgModule({
    imports: [RouterModule.forRoot(routes)],
    exports: [RouterModule]
})

export class AppRoutingModule {

}


import {NgModule} from "@angular/core";
import {RouterModule, Routes} from "@angular/router";
import {CtemplateComponent} from "ctemplate/ctemplate.component";
import {MainPageComponent} from "main-page/main-page.component";
import {AuditComponent} from "audit/audit.component";
import {ListPropertiesComponent} from "./property/list-properties.component";

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
    }
];

@NgModule({
    imports: [RouterModule.forRoot(routes)],
    exports: [RouterModule]
})

export class AppRoutingModule {

}


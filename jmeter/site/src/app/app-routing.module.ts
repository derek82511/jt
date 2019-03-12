import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { ScenarioManagerComponent } from './components/scenarioManager/scenarioManager.component';
import { JobManagerComponent } from './components/jobManager/jobManager.component';
import { JobConsoleComponent } from './components/jobConsole/jobConsole.component';

const routes: Routes = [
    { path: '', redirectTo: '/scenario', pathMatch: 'full' },
    { path: 'scenario', component: ScenarioManagerComponent },
    { path: 'job', component: JobManagerComponent },
    { path: 'jobconsole/:id', component: JobConsoleComponent, outlet: "blank" },
    { path: '**', redirectTo: '/scenario' }
];

@NgModule({
    imports: [RouterModule.forRoot(routes)],
    exports: [RouterModule]
})
export class AppRoutingModule { }

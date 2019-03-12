import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { HttpModule } from '@angular/http';
import { HttpClient, HttpClientModule, HttpClientJsonpModule, HTTP_INTERCEPTORS } from '@angular/common/http';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';

import { GridModule } from '@progress/kendo-angular-grid';
import { DateInputsModule } from '@progress/kendo-angular-dateinputs';
import { UploadModule } from '@progress/kendo-angular-upload';
import { DialogModule } from '@progress/kendo-angular-dialog';
import { LayoutModule } from '@progress/kendo-angular-layout';

import { JsonPipe } from './pipes/utilPipe';
import { CommaPipe } from './pipes/utilPipe';

import { ScenarioService } from './services/ScenarioService';
import { JobService } from './services/JobService';

import { AppComponent } from './app.component';
import { AppRoutingModule } from './app-routing.module';
import { ScenarioManagerComponent } from './components/scenarioManager/scenarioManager.component';
import { ScenarioEditFormComponent } from './components/scenarioEditForm/scenarioEditForm.component';
import { ScenarioJobEditFormComponent } from './components/scenarioJobEditForm/scenarioJobEditForm.component';
import { ScenarioJobAlertComponent } from './components/scenarioJobAlert/scenarioJobAlert.component';
import { JobManagerComponent } from './components/jobManager/jobManager.component';
import { JobConsoleComponent } from './components/jobConsole/jobConsole.component';
import { MenuComponent } from './components/menu/menu.component';

@NgModule({
  declarations: [
    AppComponent,
    ScenarioManagerComponent,
    ScenarioEditFormComponent,
    ScenarioJobEditFormComponent,
    ScenarioJobAlertComponent,
    JobManagerComponent,
    JobConsoleComponent,
    MenuComponent,
    JsonPipe,
    CommaPipe,
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    HttpModule,
    HttpClientModule,
    FormsModule,
    GridModule,
    UploadModule,
    BrowserAnimationsModule,
    DialogModule,
    DateInputsModule,
    ReactiveFormsModule,
    LayoutModule
  ],
  providers: [ScenarioService, JobService],
  bootstrap: [AppComponent]
})
export class AppModule { }

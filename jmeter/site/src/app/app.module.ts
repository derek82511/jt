import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';

import { JsonPipe, CommaPipe } from './pipes/utilPipe';

import { ScenarioService } from './services/ScenarioService';
import { JobService } from './services/JobService';

@NgModule({
  imports: [
    FormsModule,
    ReactiveFormsModule,
    JsonPipe,
    CommaPipe
  ],
  exports: [
    FormsModule,
    ReactiveFormsModule,
    JsonPipe,
    CommaPipe
  ],
  providers: [ScenarioService, JobService],
})
export class AppModule { }

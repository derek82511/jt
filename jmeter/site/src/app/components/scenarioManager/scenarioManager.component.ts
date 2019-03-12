import { Component, OnInit, Inject } from '@angular/core';
import { State } from '@progress/kendo-data-query';

import { Scenario } from '../../models/Scenario';

import { ScenarioService } from '../../services/ScenarioService';
import { JobService } from '../../services/JobService';

@Component({
    selector: 'app-scenarioManager',
    templateUrl: './scenarioManager.component.html',
    styleUrls: ['./scenarioManager.component.css']
})
export class ScenarioManagerComponent implements OnInit {
    public gridState: State = {
        skip: 0,
        take: 10
    };

    public active = false;
    public jobActive = false;
    public alertActive = false;

    public scenarios: Scenario[];

    public scenarioId: string;

    public jobId: string;

    constructor(private scenarioService: ScenarioService, private jobService: JobService) {

    }

    ngOnInit() {
        this.getScenarios();
    }

    getScenarios(): void {
        this.scenarioService.getScenarios()
            .subscribe(scenarios => this.scenarios = scenarios);
    }

    onStateChange(state: State) {
        this.gridState = state;
    }

    addHandler() {
        this.active = true;
    }

    removeHandler({ dataItem }) {
        this.scenarioService.deleteScenario(dataItem.id)
            .subscribe(() => this.getScenarios());
    }

    // editForm
    saveHandler(scenario: any) {
        this.scenarioService.createScenario(scenario)
            .subscribe(() => {
                this.getScenarios();
                this.active = false;
            });
    }

    cancelHandler() {
        this.active = false;
    }

    // jobEditForm
    addJobHandler(dataItem: any) {
        this.scenarioId = dataItem.id;
        this.jobActive = true;
    }

    saveJobHandler(job: any) {
        let params = {
            scenarioId: job.scenarioId,
            config: '',
            maxHeap: job.maxHeap,
            minHeap: job.minHeap,
            executeType: parseInt(job.executeType),
            remoteHost: '',
        };

        if (job.configsForm.length !== 0) {
            let configData = {};

            for (let i = 0; i < job.configsForm.length; i++) {
                configData[job.configsForm[i].key] = job.configsForm[i].value;
            }

            params.config = JSON.stringify(configData);
        }

        for (let i = 0; i < job.remoteHostsForm.length; i++) {
            params.remoteHost += job.remoteHostsForm[i].remoteHost;
            if (i !== job.remoteHostsForm.length - 1)
                params.remoteHost += ',';
        }

        this.jobService.createJob(params)
            .subscribe((result) => {

                //process
                this.scenarioId = null;
                this.jobActive = false;

                this.jobId = result.id;
                this.alertActive = true;
            });
    }

    cancelJobHandler() {
        this.scenarioId = null;
        this.jobActive = false;
    }

    cancelAlertHandler() {
        this.jobId = null;
        this.alertActive = false;
    }

}

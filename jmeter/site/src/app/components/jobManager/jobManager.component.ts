import { Component, OnInit } from '@angular/core';
import { State } from '@progress/kendo-data-query';

import { Constant } from '../../app.constant';

import { Job } from '../../models/Job';

import { JobService } from '../../services/JobService';

@Component({
    selector: 'app-jobManager',
    templateUrl: './jobManager.component.html',
    styleUrls: ['./jobManager.component.css']
})
export class JobManagerComponent implements OnInit {
    public gridState: State = {
        skip: 0,
        take: 10
    };

    public jobs: Job[];

    constructor(private jobService: JobService) {

    }

    ngOnInit() {
        this.getJobs();
    }

    getJobs(): void {
        this.jobService.getJobs()
            .subscribe(jobs => this.jobs = jobs);
    }

    onStateChange(state: State) {
        this.gridState = state;
    }

    report(e, dataItem: any): void {
        e.preventDefault();

        window.open(Constant.Host + dataItem.reportPath, '_blank');
    }

    jobConsole(dataItem: any) {
        window.open(window.location.origin + "/(blank:jobconsole/" + dataItem.id + ")", '_blank');
    }

}

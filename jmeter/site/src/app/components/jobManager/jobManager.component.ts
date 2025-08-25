import { Component, OnInit, ViewChild, AfterViewInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatTableDataSource, MatTableModule } from '@angular/material/table';
import { MatPaginator, MatPaginatorModule } from '@angular/material/paginator';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { CdkTableModule } from '@angular/cdk/table';

import { Constant } from '../../app.constant';
import { Job } from '../../models/Job';
import { JobService } from '../../services/JobService';
import { JsonPipe, CommaPipe } from '../../pipes/utilPipe';

@Component({
    standalone: true,
    imports: [
        CommonModule,
        MatTableModule,
        MatPaginatorModule,
        MatButtonModule,
        MatIconModule,
        CdkTableModule,
        JsonPipe,
        CommaPipe
    ],
    selector: 'app-jobManager',
    templateUrl: './jobManager.component.html',
    styleUrls: ['./jobManager.component.css']
})
export class JobManagerComponent implements OnInit, AfterViewInit {
    displayedColumns: string[] = ['scenarioName', 'createTime', 'report', 'heapSize', 'properties', 'remoteHost', 'status', 'actions'];
    dataSource = new MatTableDataSource<Job>();

    @ViewChild(MatPaginator) paginator!: MatPaginator;

    constructor(private jobService: JobService) { }

    ngOnInit() {
        this.getJobs();
    }

    ngAfterViewInit() {
        this.dataSource.paginator = this.paginator;
    }

    getJobs(): void {
        this.jobService.getJobs()
            .subscribe(jobs => {
                this.dataSource.data = jobs;
            });
    }

    report(e: Event, dataItem: Job): void {
        e.preventDefault();
        window.open(Constant.Host + dataItem.reportPath, '_blank');
    }

    jobConsole(dataItem: Job) {
        window.open(window.location.origin + "/(blank:jobconsole/" + dataItem.id + ")", '_blank');
    }
}

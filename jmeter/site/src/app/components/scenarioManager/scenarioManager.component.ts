import { Component, OnInit, ViewChild, AfterViewInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatDialog, MatDialogModule } from '@angular/material/dialog';
import { MatTableDataSource, MatTableModule } from '@angular/material/table';
import { MatPaginator, MatPaginatorModule } from '@angular/material/paginator';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatToolbarModule } from '@angular/material/toolbar';

import { Scenario } from '../../models/Scenario';
import { ScenarioService } from '../../services/ScenarioService';
import { JobService } from '../../services/JobService';
import { ScenarioEditFormComponent } from '../scenarioEditForm/scenarioEditForm.component';
import { ScenarioJobEditFormComponent } from '../scenarioJobEditForm/scenarioJobEditForm.component';
import { ScenarioJobAlertComponent } from '../scenarioJobAlert/scenarioJobAlert.component';

@Component({
    standalone: true,
    imports: [
        CommonModule,
        MatDialogModule,
        MatTableModule,
        MatPaginatorModule,
        MatButtonModule,
        MatIconModule,
        MatToolbarModule
    ],
    selector: 'app-scenarioManager',
    templateUrl: './scenarioManager.component.html',
    styleUrls: ['./scenarioManager.component.css']
})
export class ScenarioManagerComponent implements OnInit, AfterViewInit {
    displayedColumns: string[] = ['name', 'scriptName', 'createTime', 'actions'];
    dataSource = new MatTableDataSource<Scenario>();

    @ViewChild(MatPaginator) paginator!: MatPaginator;

    constructor(
        private scenarioService: ScenarioService,
        private jobService: JobService,
        public dialog: MatDialog
    ) { }

    ngOnInit() {
        this.getScenarios();
    }

    ngAfterViewInit() {
        this.dataSource.paginator = this.paginator;
    }

    getScenarios(): void {
        this.scenarioService.getScenarios()
            .subscribe(scenarios => {
                this.dataSource.data = scenarios;
            });
    }

    addHandler() {
        const dialogRef = this.dialog.open(ScenarioEditFormComponent, {
            width: '400px'
        });

        dialogRef.afterClosed().subscribe(result => {
            if (result) {
                this.saveHandler(result);
            }
        });
    }

    removeHandler(scenario: Scenario) {
        this.scenarioService.deleteScenario(scenario.id)
            .subscribe(() => this.getScenarios());
    }

    saveHandler(scenario: any) {
        this.scenarioService.createScenario(scenario)
            .subscribe(() => {
                this.getScenarios();
            });
    }

    addJobHandler(dataItem: any) {
        const dialogRef = this.dialog.open(ScenarioJobEditFormComponent, {
            width: '600px',
            data: { scenarioId: dataItem.id }
        });

        dialogRef.afterClosed().subscribe(result => {
            if (result) {
                this.saveJobHandler(result);
            }
        });
    }

    saveJobHandler(job: any) {
        let configData: { [key: string]: string } = {};
        if (job.config && Array.isArray(job.config)) {
            job.config.forEach((item: { key: string, value: string }) => {
                if (item.key) {
                    configData[item.key] = item.value;
                }
            });
        }

        let remoteHostData = '';
        if (job.remoteHost && Array.isArray(job.remoteHost)) {
            remoteHostData = job.remoteHost
                .map((item: { remoteHost: string }) => item.remoteHost)
                .filter((host: string | null) => host)
                .join(',');
        }

        let params = {
            scenarioId: job.scenarioId,
            config: configData,
            maxHeap: job.maxHeap,
            minHeap: job.minHeap,
            executeType: parseInt(job.executeType, 10),
            remoteHost: remoteHostData,
        };

        this.jobService.createJob(params)
            .subscribe((result) => {
                if (result) {
                    this.dialog.open(ScenarioJobAlertComponent, {
                        data: { jobId: result.id }
                    });
                }
            });
    }
}

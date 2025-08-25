import { Component, OnInit, ElementRef, ViewChild } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { CommonModule } from '@angular/common';
import { MatDialog, MatDialogModule } from '@angular/material/dialog';
import { MatButtonModule } from '@angular/material/button';
import { MatTabsModule } from '@angular/material/tabs';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';

import { Constant } from '../../app.constant';
import { Job } from '../../models/Job';
import { JobService } from '../../services/JobService';
import { ScenarioJobAlertComponent } from '../scenarioJobAlert/scenarioJobAlert.component';
import { JsonPipe, CommaPipe } from '../../pipes/utilPipe';

@Component({
    standalone: true,
    imports: [
        CommonModule,
        MatDialogModule,
        MatButtonModule,
        MatTabsModule,
        MatProgressSpinnerModule,
        JsonPipe,
        CommaPipe
    ],
    selector: 'app-jobConsole',
    templateUrl: './jobConsole.component.html',
    styleUrls: ['./jobConsole.component.css']
})
export class JobConsoleComponent implements OnInit {
    private socket: WebSocket | null = null;

    public job!: Job;

    public lines: string[] = [];

    public firstStartup = false;

    @ViewChild('scrollContainer') private scrollContainer!: ElementRef;

    constructor(
        private activatedRoute: ActivatedRoute,
        private jobService: JobService,
        public dialog: MatDialog
    ) { }

    ngOnInit() {
        const id = this.activatedRoute.snapshot.paramMap.get('id');
        if (id) {
            this.getJob(id);
        }
    }

    getJob(id: string): void {
        this.jobService.getJob(id).subscribe((job) => {
            if (job) {
                this.job = job;
            }
        });
    }

    run(): void {
        this.socket = new WebSocket(this.getWsUrl());

        this.socket.onopen = (event) => {
            console.log("Status: Connected");
            const registerMsg = { event: "register", payload: this.job.id };
            this.socket?.send(JSON.stringify(registerMsg));
        };

        this.socket.onclose = (event) => {
            console.log("Status: Disconnected");
        };

        this.socket.onerror = (error) => {
            console.error("WebSocket Error:", error);
        };

        this.socket.onmessage = (event) => {
            const msg = JSON.parse(event.data);

            switch (msg.event) {
                case "onRegister":
                    if (msg.payload === 'ok') {
                        this.jobService.runJob(this.job.id).subscribe(() => {
                            this.job.status = "1";
                        });
                    }
                    break;
                case "console":
                    this.lines.push(msg.payload);
                    this.scrollToBottom();
                    break;
                case "finish":
                    if (msg.payload === 'ok') {
                        this.job.status = "2";
                        this.firstStartup = true;
                    }
                    break;
            }
        };
    }

    terminate(): void {
        this.jobService.terminateJob(this.job.id).subscribe((result: { recovery: boolean }) => {
            if (result.recovery) {
                this.job.status = "2";
                this.firstStartup = true;
            }
        });
        this.socket?.close();
    }

    restart(): void {
        let configData = {};
        try {
            // The config from the job object is a JSON string; parse it back to an object.
            if (this.job.config) {
                configData = JSON.parse(this.job.config);
            }
        } catch (e) {
            console.error('Could not parse job config JSON:', e);
            // Optionally, show an error to the user.
        }

        const params = {
            scenarioId: this.job.scenarioId,
            config: configData,
            maxHeap: this.job.maxHeap,
            minHeap: this.job.minHeap,
            executeType: parseInt(this.job.executeType, 10),
            remoteHost: this.job.remoteHost,
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

    report(): void {
        window.open(Constant.Host + this.job.reportPath, '_blank');
    }

    private scrollToBottom(): void {
        setTimeout(() => {
            try {
                this.scrollContainer.nativeElement.scrollTop = this.scrollContainer.nativeElement.scrollHeight;
            } catch (err) { }
        }, 100);
    }

    private getWsUrl(): string {
        const wsHost = Constant.Host.replace("http", "ws");
        return `${wsHost}/ws`;
    }
}

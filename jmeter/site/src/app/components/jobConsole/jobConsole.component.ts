import { Component, OnInit, ElementRef, ViewChild } from '@angular/core';
import { ActivatedRoute } from '@angular/router';

import { Constant } from '../../app.constant';

import { Job } from '../../models/Job';

import { JobService } from '../../services/JobService';

@Component({
    selector: 'app-jobConsole',
    templateUrl: './jobConsole.component.html',
    styleUrls: ['./jobConsole.component.css']
})
export class JobConsoleComponent implements OnInit {
    private socket: any;

    public job: Job;

    public lines = [];

    public firstStartup = false;

    public alertActive = false;
    public alertJobId: string;

    @ViewChild('scrollContainer') private scrollContainer: ElementRef;

    constructor(private activatedRoute: ActivatedRoute, private jobService: JobService) {

    }

    ngOnInit() {
        var id = this.activatedRoute.snapshot.paramMap.get('id');

        this.getJob(id);
    }

    getJob(id: string): void {
        this.jobService.getJob(id).subscribe((job) => {
            this.job = job;
        });
    }

    run(): void {
        this.socket = new window['Ws'](this.getWsUrl())

        this.socket.OnConnect(() => {
            console.log("Status: Connected");

            this.socket.Emit('register', this.job.id);
        });

        this.socket.OnDisconnect(() => {
            console.log("Status: Disconnected");
        });

        this.socket.On("onRegister", (msg) => {
            if (msg === 'ok') {
                this.jobService.runJob(this.job.id).subscribe(() => {
                    this.job.status = "1";
                });
            }
        });

        this.socket.On("console", (consoleLine) => {
            this.lines.push(consoleLine)
            setTimeout(() => {
                try {
                    this.scrollContainer.nativeElement.scrollTop = this.scrollContainer.nativeElement.scrollHeight;
                } catch (err) { }
            })
        });

        this.socket.On("finish", (msg) => {
            if (msg === 'ok') {
                this.job.status = "2";
                this.firstStartup = true;
            }
        });
    }

    terminate(): void {
        this.jobService.terminateJob(this.job.id).subscribe((result) => {
            if (result.recovery === "1") {
                this.job.status = "2";
                this.firstStartup = true;
            }
        });
    }

    restart(): void {
        let params = {
            scenarioId: this.job.scenarioId,
            config: this.job.config,
            maxHeap: this.job.maxHeap,
            minHeap: this.job.minHeap,
            executeType: parseInt(this.job.executeType),
            remoteHost: this.job.remoteHost,
        };

        this.jobService.createJob(params)
            .subscribe((result) => {
                this.alertJobId = result.id;
                this.alertActive = true;
            });
    }

    cancelAlertHandler() {
        this.alertJobId = null;
        this.alertActive = false;
    }

    report(): void {
        window.open(Constant.Host + this.job.reportPath, '_blank');
    }

    private getWsUrl(): string {
        if (Constant.Host === "") {
            var loc = window.location, new_uri;
            if (loc.protocol === "https:") {
                new_uri = "wss:";
            } else {
                new_uri = "ws:";
            }
            new_uri += "//" + loc.host;
            new_uri += "/echo";

            return new_uri
        } else {
            return Constant.Host.replace("http", "ws") + "/echo"
        }
    }

}

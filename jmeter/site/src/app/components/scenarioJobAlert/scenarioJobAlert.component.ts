import { Component, Inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatDialogModule, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { MatButtonModule } from '@angular/material/button';

@Component({
    standalone: true,
    imports: [
        CommonModule,
        MatDialogModule,
        MatButtonModule
    ],
    selector: 'app-scenarioJobAlert',
    templateUrl: './scenarioJobAlert.component.html',
    styleUrls: ['./scenarioJobAlert.component.css']
})
export class ScenarioJobAlertComponent {
    constructor(@Inject(MAT_DIALOG_DATA) public data: { jobId: string }) {

    }

    public jobConsole(e: Event): void {
        e.preventDefault();

        window.open(window.location.origin + "/(blank:jobconsole/" + this.data.jobId + ")", '_blank');
    }
}

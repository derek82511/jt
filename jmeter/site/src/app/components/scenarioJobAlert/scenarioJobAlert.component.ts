import { Component, Input, Output, EventEmitter } from '@angular/core';

import { Constant } from '../../app.constant';

@Component({
    selector: 'app-scenarioJobAlert',
    templateUrl: './scenarioJobAlert.component.html',
    styleUrls: ['./scenarioJobAlert.component.css']
})
export class ScenarioJobAlertComponent {
    @Input() jobId: string;

    @Output() cancel: EventEmitter<any> = new EventEmitter();

    constructor() {

    }

    public jobConsole(e): void {
        e.preventDefault();

        window.open(window.location.origin + "/(blank:jobconsole/" + this.jobId + ")", '_blank');
    }

    public closeForm(): void {
        this.cancel.emit();
    }
}

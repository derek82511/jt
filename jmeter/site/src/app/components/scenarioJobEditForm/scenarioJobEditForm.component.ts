import { Component, Inject } from '@angular/core';
import { Validators, FormGroup, FormControl, FormArray, ReactiveFormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common';
import { MatDialogModule, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { MatButtonModule } from '@angular/material/button';
import { MatInputModule } from '@angular/material/input';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatRadioModule } from '@angular/material/radio';
import { MatTabsModule } from '@angular/material/tabs';
import { MatIconModule } from '@angular/material/icon';

@Component({
    standalone: true,
    imports: [
        CommonModule,
        ReactiveFormsModule,
        MatDialogModule,
        MatButtonModule,
        MatFormFieldModule,
        MatInputModule,
        MatRadioModule,
        MatTabsModule,
        MatIconModule
    ],
    selector: 'app-scenarioJobEditForm',
    templateUrl: './scenarioJobEditForm.component.html',
    styleUrls: ['./scenarioJobEditForm.component.css']
})
export class ScenarioJobEditFormComponent {
    public editForm: FormGroup;
    
    constructor(@Inject(MAT_DIALOG_DATA) public data: { scenarioId: string }) {
        const remoteHost = localStorage.getItem('remoteHost');
        let initialHosts: FormGroup[] = [];
        let executeType = '0';

        if (remoteHost) {
            const hosts = remoteHost.split(';').filter(h => h);
            if (hosts.length > 0) {
                executeType = '1';
                initialHosts = hosts.map(host => new FormGroup({ 'remoteHost': new FormControl(host) }));
            }
        }

        if (initialHosts.length === 0) {
            initialHosts.push(new FormGroup({ 'remoteHost': new FormControl('') }));
        }
        
        this.editForm = new FormGroup({
            'minHeap': new FormControl('', Validators.required),
            'maxHeap': new FormControl('', Validators.required),
            'configs': new FormArray([]),
            'executeType': new FormControl(executeType),
            'remoteHosts': new FormArray(initialHosts),
        });
    }

    get configs(): FormArray {
        return this.editForm.get('configs') as FormArray;
    }

    addConfig(): void {
        this.configs.push(new FormGroup({
            'key': new FormControl('', Validators.required),
            'value': new FormControl('', Validators.required),
        }));
    }

    deleteConfig(index: number): void {
        this.configs.removeAt(index);
    }
    
    get remoteHosts(): FormArray {
        return this.editForm.get('remoteHosts') as FormArray;
    }
    
    addRemoteHost(): void {
        this.remoteHosts.push(new FormGroup({
            'remoteHost': new FormControl(''),
        }));
    }

    deleteRemoteHost(index: number): void {
        this.remoteHosts.removeAt(index);
    }
    
    public getSaveData() {
        if (!this.editForm.valid) {
            return null;
        }

        const formValue = this.editForm.value;
        const result: any = {
            scenarioId: this.data.scenarioId,
            minHeap: formValue.minHeap,
            maxHeap: formValue.maxHeap,
            executeType: formValue.executeType,
            config: {},
            remoteHost: ''
        };

        result.config = formValue.configs.reduce((acc: any, curr: any) => {
            acc[curr.key] = curr.value;
            return acc;
        }, {});
        
        if (formValue.executeType === '1') {
            result.remoteHost = formValue.remoteHosts
                .map((host: any) => host.remoteHost)
                .filter((host: any) => host)
                .join(',');
        }

        return result;
    }
}

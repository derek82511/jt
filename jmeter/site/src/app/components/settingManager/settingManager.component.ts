import { Component, OnInit } from '@angular/core';
import { FormGroup, FormControl, FormArray, ReactiveFormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatRadioModule } from '@angular/material/radio';

@Component({
    standalone: true,
    imports: [
        CommonModule,
        ReactiveFormsModule,
        MatButtonModule,
        MatIconModule,
        MatFormFieldModule,
        MatInputModule,
        MatRadioModule
    ],
    selector: 'app-settingManager',
    templateUrl: './settingManager.component.html',
    styleUrls: ['./settingManager.component.css']
})
export class SettingManagerComponent implements OnInit {
    public editForm!: FormGroup;

    ngOnInit() {
        const remoteHost = localStorage.getItem('remoteHost');
        let initialHosts: FormGroup[] = [];
        let executeType = '0';

        if (remoteHost) {
            const hosts = remoteHost.split(';').filter(h => h);
            if (hosts.length > 0) {
                executeType = '1';
                initialHosts = hosts.map(host => new FormGroup({
                    'remoteHost': new FormControl(host),
                }));
            }
        }
        
        if (initialHosts.length === 0) {
            initialHosts.push(new FormGroup({ 'remoteHost': new FormControl('') }));
        }

        this.editForm = new FormGroup({
            'executeType': new FormControl(executeType),
            'remoteHosts': new FormArray(initialHosts),
        });
    }

    get remoteHosts(): FormArray {
        return this.editForm.get('remoteHosts') as FormArray;
    }

    public onSave(): void {
        if (this.editForm.value['executeType'] === '0') {
            localStorage.setItem('remoteHost', '');
        } else {
            const remoteHost = this.editForm.value['remoteHosts']
                .map((host: { remoteHost: string }) => host.remoteHost)
                .filter((host: string) => host)
                .join(';');
            localStorage.setItem('remoteHost', remoteHost);
        }
    }

    public addRemoteHost(): void {
        this.remoteHosts.push(new FormGroup({
            'remoteHost': new FormControl(''),
        }));
    }

    public deleteRemoteHost(index: number): void {
        this.remoteHosts.removeAt(index);
    }
}

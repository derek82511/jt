import { Component, OnInit, Inject, Input, Output, EventEmitter } from '@angular/core';
import { Validators, FormGroup, FormControl, FormArray } from '@angular/forms';

@Component({
    selector: 'app-settingManager',
    templateUrl: './settingManager.component.html',
    styleUrls: ['./settingManager.component.css']
})
export class SettingManagerComponent implements OnInit {
    public configsForm = new FormArray([]);

    public remoteHostsForm: FormArray;

    public editForm: FormGroup;

    ngOnInit() {
        let remoteHost = localStorage.getItem('remoteHost');

        if (remoteHost) {
            let hosts = remoteHost.split(';');

            this.remoteHostsForm = new FormArray([]);

            hosts.forEach(host => {
                if (host) {
                    this.remoteHostsForm.push(new FormGroup({
                        'remoteHost': new FormControl(host),
                    }));
                }
            });

            this.editForm = new FormGroup({
                'executeType': new FormControl('1'),
                'remoteHostsForm': this.remoteHostsForm,
            });
        } else {
            this.remoteHostsForm = new FormArray([new FormGroup({
                'remoteHost': new FormControl(''),
            })]);

            this.editForm = new FormGroup({
                'executeType': new FormControl('0'),
                'remoteHostsForm': this.remoteHostsForm,
            });
        }
    }

    public onSave(e): void {
        e.preventDefault();

        if (this.editForm.value['executeType'] === '0') {
            this.editForm.value['remoteHostsForm'] = [];
        } else {
            for (let i = 0; i < this.editForm.value['remoteHostsForm'].length; i++) {
                if (this.editForm.value['remoteHostsForm'][i].remoteHost === '') {
                    //
                    return;
                }
            }
        }

        // this.save.emit(this.editForm.value);

        let remoteHost = '';

        this.editForm.value['remoteHostsForm'].forEach(element => {
            remoteHost += element.remoteHost + ';';
        });

        localStorage.setItem('remoteHost', remoteHost);
    }

    public addRemoteHost(e): void {
        e.preventDefault();

        this.remoteHostsForm.push(new FormGroup({
            'remoteHost': new FormControl(''),
        }));
    }
    public deleteRemoteHost(e, index): void {
        e.preventDefault();

        this.remoteHostsForm.removeAt(index);
    }
}

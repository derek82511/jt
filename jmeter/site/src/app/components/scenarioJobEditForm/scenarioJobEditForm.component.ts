import { Component, Input, Output, EventEmitter } from '@angular/core';
import { Validators, FormGroup, FormControl, FormArray } from '@angular/forms';

@Component({
    selector: 'app-scenarioJobEditForm',
    templateUrl: './scenarioJobEditForm.component.html',
    styleUrls: ['./scenarioJobEditForm.component.css']
})
export class ScenarioJobEditFormComponent {
    public configsForm = new FormArray([]);

    public remoteHostsForm = new FormArray([new FormGroup({
        'remoteHost': new FormControl(''),
    })]);

    public editForm: FormGroup = new FormGroup({
        'minHeap': new FormControl('', Validators.required),
        'maxHeap': new FormControl('', Validators.required),
        'config': new FormControl('{}'),
        'configsForm': this.configsForm,
        'executeType': new FormControl('0'),
        'remoteHostsForm': this.remoteHostsForm,
    });

    public propertyType: number = 0;

    @Input() scenarioId: string;

    @Output() cancel: EventEmitter<any> = new EventEmitter();
    @Output() save: EventEmitter<any> = new EventEmitter();

    constructor() {

    }

    public onSave(e): void {
        e.preventDefault();

        if (this.propertyType !== 0) {
            this.changePropertyType(e, 0);
        }

        try {
            this.editForm.value['config'] = JSON.parse(this.editForm.value['config']);
        } catch (e) {
            this.editForm.value['config'] = {};
        }

        this.configsForm = new FormArray([]);

        for (let key in this.editForm.value['config']) {
            this.configsForm.push(new FormGroup({
                'key': new FormControl(key),
                'value': new FormControl(this.editForm.value['config'][key])
            }));
        }

        this.editForm.setControl('configsForm', this.configsForm);

        this.editForm.value['scenarioId'] = this.scenarioId;

        for (let i = 0; i < this.editForm.value['configsForm'].length; i++) {
            if (this.editForm.value['configsForm'][i].key === '' || this.editForm.value['configsForm'][i].value === '') {
                //
                return;
            }
        }

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

        this.save.emit(this.editForm.value);
    }

    public onCancel(e): void {
        e.preventDefault();

        this.closeForm();
    }

    public closeForm(): void {
        this.cancel.emit();
    }

    public changePropertyType(e, type): void {
        e.preventDefault();

        this.propertyType = type;

        if (this.propertyType === 0) {
            let configData = {};

            for (let i = 0; i < this.editForm.value['configsForm'].length; i++) {
                configData[this.editForm.value['configsForm'][i].key] = this.editForm.value['configsForm'][i].value;
            }

            this.editForm.value['config'] = JSON.stringify(configData, null, 4);

            this.editForm.reset(this.editForm.value);
        } else {
            try {
                this.editForm.value['config'] = JSON.parse(this.editForm.value['config']);
            } catch (e) {
                this.editForm.value['config'] = {};
            }

            this.configsForm = new FormArray([]);

            for (let key in this.editForm.value['config']) {
                this.configsForm.push(new FormGroup({
                    'key': new FormControl(key),
                    'value': new FormControl(this.editForm.value['config'][key])
                }));
            }

            this.editForm.setControl('configsForm', this.configsForm);
        }
    }

    public addConfig(e): void {
        e.preventDefault();

        this.configsForm.push(new FormGroup({
            'key': new FormControl(''),
            'value': new FormControl(''),
        }));
    }
    public deleteConfig(e, index): void {
        e.preventDefault();

        this.configsForm.removeAt(index);
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

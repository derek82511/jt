import { Component, Output, EventEmitter } from '@angular/core';
import { Validators, FormGroup, FormControl } from '@angular/forms';

@Component({
    selector: 'app-scenarioEditForm',
    templateUrl: './scenarioEditForm.component.html',
    styleUrls: ['./scenarioEditForm.component.css']
})
export class ScenarioEditFormComponent {
    public files;

    public editForm: FormGroup = new FormGroup({
        'name': new FormControl('', Validators.required),
        'uploadedFile': new FormControl()
    });

    @Output() cancel: EventEmitter<any> = new EventEmitter();
    @Output() save: EventEmitter<any> = new EventEmitter();

    constructor() {

    }

    onFileChange(event) {
        if (event.target.files && event.target.files.length >= 1) {
            this.files = event.target.files;
        }
    }

    public onSave(e): void {
        e.preventDefault();

        if (!this.files) {
            return
        }

        this.editForm.value['uploadedFile'] = this.files;

        this.save.emit(this.editForm.value);
    }

    public onCancel(e): void {
        e.preventDefault();

        this.closeForm();
    }

    public closeForm(): void {
        this.cancel.emit();
    }
}

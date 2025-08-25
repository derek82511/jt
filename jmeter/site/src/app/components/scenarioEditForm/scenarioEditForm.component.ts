import { Component } from '@angular/core';
import { Validators, FormGroup, FormControl, ReactiveFormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common';
import { MatDialogModule } from '@angular/material/dialog';
import { MatButtonModule } from '@angular/material/button';
import { MatInputModule } from '@angular/material/input';
import { MatFormFieldModule } from '@angular/material/form-field';

@Component({
    standalone: true,
    imports: [
        CommonModule,
        ReactiveFormsModule,
        MatDialogModule,
        MatButtonModule,
        MatFormFieldModule,
        MatInputModule
    ],
    selector: 'app-scenarioEditForm',
    templateUrl: './scenarioEditForm.component.html',
    styleUrls: ['./scenarioEditForm.component.css']
})
export class ScenarioEditFormComponent {
    public files!: FileList;

    public editForm: FormGroup = new FormGroup({
        'name': new FormControl('', Validators.required),
        'uploadedFile': new FormControl(null, Validators.required)
    });

    constructor() {

    }

    onFileChange(event: Event) {
        const input = event.target as HTMLInputElement;
        if (input.files && input.files.length > 0) {
            this.files = input.files;
            this.editForm.patchValue({
                uploadedFile: this.files
            });
        }
    }

    public getSaveData() {
        if (this.editForm.valid) {
            return this.editForm.value;
        }
        return null;
    }
}

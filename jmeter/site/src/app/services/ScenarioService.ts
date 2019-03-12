import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';

import { Observable } from 'rxjs/Observable';
import { of } from 'rxjs/observable/of';
import { catchError, map, tap } from 'rxjs/operators';

import { Constant } from '../app.constant';
import { Scenario } from '../models/Scenario';

@Injectable()
export class ScenarioService {
    private url = Constant.Host + "/api/scenario"

    constructor(private http: HttpClient) { }

    getScenarios(): Observable<Scenario[]> {
        return this.http.get<Scenario[]>(this.url)
            .pipe(catchError(this.handleError('getScenarios', [])));
    }

    createScenario(data: any): Observable<any> {
        const formData: FormData = new FormData();
        formData.append("name", data.name);
        for (let i = 0; i < data.uploadedFile.length; i++) {
            formData.append("uploadedFile", data.uploadedFile[i]);
        }
        return this.http.post(this.url, formData)
            .pipe(catchError(this.handleError('createScenario', null)));
    }

    deleteScenario(id: string): Observable<any> {
        return this.http.delete(this.url + "/" + id)
            .pipe(catchError(this.handleError('deleteScenario', null)));
    }

    private handleError<T>(operation = 'operation', result?: T) {
        return (error: any): Observable<T> => {
            console.error(error);
            return of(result as T);
        };
    }
}

import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';

import { Observable } from 'rxjs/Observable';
import { of } from 'rxjs/observable/of';
import { catchError } from 'rxjs/operators';

import { Constant } from '../app.constant';

import { Job } from '../models/Job';

@Injectable()
export class JobService {
    private url = Constant.Host + "/api/job"

    constructor(private http: HttpClient) { }

    getJobs(): Observable<Job[]> {
        return this.http.get<Job[]>(this.url)
            .pipe(catchError(this.handleError('getJobs', [])));
    }

    getJob(id: string): Observable<Job> {
        return this.http.get<Job>(this.url + "/" + id)
            .pipe(catchError(this.handleError('getJob', null)));
    }

    createJob(data: any): Observable<Job> {
        return this.http.post(this.url, data)
            .pipe(catchError(this.handleError('createJob', null)));
    }

    runJob(id: string): Observable<any> {
        return this.http.post(this.url + "/run/" + id, "")
            .pipe(catchError(this.handleError('runJob', null)));
    }

    terminateJob(id: string): Observable<any> {
        return this.http.post(this.url + "/terminate/" + id, "")
            .pipe(catchError(this.handleError('terminateJob', null)));
    }

    private handleError<T>(operation = 'operation', result?: T) {
        return (error: any): Observable<T> => {
            console.error(error);
            return of(result as T);
        };
    }
}

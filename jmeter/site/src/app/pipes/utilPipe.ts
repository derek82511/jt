import { Pipe, PipeTransform } from '@angular/core';

@Pipe({ name: 'jsonFilter' })
export class JsonPipe implements PipeTransform {
    transform(value: any): string {
        return JSON.stringify(JSON.parse(value), null, 4).replace(/(?:\r\n|\r|\n)/g, '<br>').replace(/ /g, '&nbsp;');
    }
}

@Pipe({ name: 'commaFilter' })
export class CommaPipe implements PipeTransform {
    transform(value: any): string {
        let out = '';

        let datas = value.split(',');
        datas.forEach(data => {
            out += data + '<br>';
        });

        return out;
    }
}

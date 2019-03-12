import { Component, OnInit } from '@angular/core';

@Component({
    selector: 'app-menu',
    templateUrl: './menu.component.html',
    styleUrls: ['./menu.component.css']
})
export class MenuComponent implements OnInit {
    public menus: any[];
    constructor() { }

    ngOnInit() {
        this.menus = [
            { url: '/scenario', name: 'Test Scenario' },
            { url: '/job', name: 'Test Job' }
        ];
    }

}

import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { MatListModule } from '@angular/material/list';
import { MatIconModule } from '@angular/material/icon';

@Component({
    standalone: true,
    imports: [
        CommonModule,
        RouterModule,
        MatListModule,
        MatIconModule
    ],
    selector: 'app-menu',
    templateUrl: './menu.component.html',
    styleUrls: ['./menu.component.css']
})
export class MenuComponent implements OnInit {
    public menus: any[] = [];
    constructor() { }

    ngOnInit() {
        this.menus = [
            { name: 'Test Scenario', url: '/scenario', icon: 'description' },
            { name: 'Test Job', url: '/job', icon: 'play_circle_filled' },
            { name: 'Setting', url: '/setting', icon: 'settings' }
        ];
    }

}

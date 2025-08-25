import { Component } from '@angular/core';
import { Location, CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';

import { Constant } from './app.constant';
import { MenuComponent } from './components/menu/menu.component';

@Component({
  standalone: true,
  imports: [
    CommonModule,
    RouterModule,
    MenuComponent
  ],
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  public version = Constant.Version;

  public isDefault = true;

  constructor(private location: Location) { }

  ngOnInit() {
    if (this.location.path().includes("(") && this.location.path().includes(")")) {
      this.isDefault = false
    }
  }
}

import { Component } from '@angular/core';
import { Location } from '@angular/common';

import { Constant } from './app.constant';

@Component({
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

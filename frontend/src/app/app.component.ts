import { Component, OnInit } from '@angular/core';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  //styles: ['styles.scss']
   styleUrls: ['./app.component.css']
})

export class AppComponent implements OnInit {
  title = 'app';
  shouldRun = true;
  // edited = true;

  _toggleAutoCollapseHeight(): void {
    
  }

  _toggleSidebar():void {
    
  }

  ngOnInit(): void {
    // let button = <HTMLElement>document.body.querySelector(".navbar-brand");
    //button.addEventListener("click", () => { alert('test')});
  }
  
  // hideNav(): void {
  //   this.edited = false;
  // }

}

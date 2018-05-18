import { Component, OnInit, Pipe, PipeTransform, Injectable } from '@angular/core';
import { DomSanitizer} from '@angular/platform-browser';

@Pipe({ name: 'safe' })
export class SafePipe implements PipeTransform {
  constructor(private sanitizer: DomSanitizer) {}
  transform(url) {
    return this.sanitizer.bypassSecurityTrustResourceUrl(url);
  }
} 


function _window() : any {
    // return the global native browser window object
    return window;
  }
  
  @Injectable()
  export class WindowRef {
    get nativeWindow() : any {
       return _window();
    }
  }

  
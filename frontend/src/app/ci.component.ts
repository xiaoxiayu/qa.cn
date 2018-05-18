import { Component, OnInit } from '@angular/core';
import { DomSanitizer } from '@angular/platform-browser';
import { WindowRef } from './fxqacommon.component';

import 'rxjs/add/operator/switchMap';
import { ActivatedRoute, ParamMap } from '@angular/router';
import { CIService } from './ci.service';
//import * as $ from 'jquery';



class CIData {
  id: string;
  url: string;
  constructor(id: string, url: string) {
    this.id = id;
    this.url = url;
  }
}


@Component({
  selector: 'app-ci',
  templateUrl: './ci.component.html',
  styleUrls: ['./ci.component.css']
})


export class CIComponent implements OnInit {
  title = 'ci';
  ci_servers: CIData[];
  frame_h: number;
  selectedIndex: number = 0;

  constructor(
    private winRef: WindowRef,
    private route: ActivatedRoute,
    private service: CIService,) {
    console.log('Window object', winRef.nativeWindow);
  }

  ngOnInit(): void {
    this.frame_h = this.winRef.nativeWindow.screen.availHeight;
    this.getCIs();
    
    this.route.paramMap
            .switchMap((params: ParamMap) =>
            this.service.updateTab(params.get('id')))
            .subscribe(select_id => this.selectedIndex = +select_id);

  //   document.addEventListener('keyup', function() {
  //     console.log('keys pressed');
  //  });

  }

  tabChanged(event: any) {
    

  }

  getCIs(): void {
    this.ci_servers = [
      new CIData('fxcore', 'http://10.103.2.160:8080/'), 
      new CIData('sdk', 'http://10.103.129.206:8080/'), 
      new CIData('set', 'http://10.103.2.202:8080/'), 
      new CIData('phantom-reader-package', 'http://10.103.2.254:8080/'), 
      new CIData('phantom-reader-test', 'http://10.103.129.28:8080/')
    ];
  }

  gotoDetail(): void {
    //this.frame_h = $(window).height();
    //alert(this.frame_h);
  }
}

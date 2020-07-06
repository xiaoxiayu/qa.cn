import { Component, OnInit, ViewChild, ElementRef } from '@angular/core';
import { DomSanitizer } from '@angular/platform-browser';
import { HttpClientModule, HttpClient, HttpHeaders } from '@angular/common/http';

//import * as $ from 'jquery';



@Component({
    selector: 'app-xxjs',
    templateUrl: './xxjs.component.html',
    styleUrls: ['./xxjs.component.css']
})


export class xxJSComponent implements OnInit {
    text: string = "";
    win_ip: string = '10.203.13.3:9091';
    getlog_obj: any;
    options: any = { maxLines: 1000, printMargin: false };
    run: boolean = false;

    @ViewChild("respondtext") respondtext: ElementRef;

    constructor(private http: HttpClient, ) {
    }

    ngOnInit(): void {


    }

    onChange(code) {
        // console.log("new code", code);
    }

    b64EncodeUnicode(str: string): string {
        return btoa(encodeURIComponent(str).replace(/%([0-9A-F]{2})/g,
            function toSolidBytes(match, p1) {
                return String.fromCharCode(+('0x' + p1));
            }));
    }

    clearOutput(): void {
        this.respondtext.nativeElement.value = '';
    }

    clearCode(): void {
        this.text = '';
    }

    get_current_time(): string {
        var a = new Date();
        var date = a.getDate();
        var hour = a.getHours() < 10 ? '0' + a.getHours() : a.getHours();
        var min = a.getMinutes() < 10 ? '0' + a.getMinutes() : a.getMinutes();
        var sec = a.getSeconds() < 10 ? '0' + a.getSeconds() : a.getSeconds();
        var time = hour + ":" + min + ":" + sec;
        return time;
    }

    getlog(): void {
        let url_s: string = 'http://' + this.win_ip + '/code';
        this.http
            .get(url_s).subscribe(resp => {
                console.info(resp);
                clearInterval(this.getlog_obj);
                // if (logres.length > 0) {
                //     $("#respond-text").append(this.get_current_time()+":" + logres + "\n");
                //     clearInterval(getlog_obj);
                // }
            });


    }


    setCode(): void {
        if (this.text == '') {
            return;
        }
        this.run = true;
        let url_s: string = 'http://' + this.win_ip + '/code';
        const body = "code=" + this.b64EncodeUnicode(this.text);
        // console.info(this.respondtext);
        console.info(url_s);
        this.http
            .post(url_s, body, {
                headers: new HttpHeaders().set('Content-Type', 'application/x-www-form-urlencoded'),
            }).subscribe(resp => {
                console.info(resp);
                if (resp['ret'] == -1) {
                    if (resp['msg'] == 'running') {

                        this.respondtext.nativeElement.value += this.get_current_time() + ": " + 'test server busy.wait and retry.\n';
                        // $("#respond-text").append(time+": "+ 'test server busy.wait and retry.\n');
                        // NProgress.done();
                    }
                } else {
                    this.getlog_obj = setInterval(() => { this.getlog(); }, 1000);
                }
                this.run = false;
            },
                errresp => {
                    this.run = false;
                    // console.info(errresp);
                    // if (errresp.status == 404) {
                    this.respondtext.nativeElement.value += this.get_current_time() + ": " + 'Error: The JS server did not start.\n';
                    return;
                    // }
                }
            );

    }

}

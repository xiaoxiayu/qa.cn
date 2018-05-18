import { Component, OnInit, ElementRef, ViewChild } from '@angular/core';
import { DomSanitizer } from '@angular/platform-browser';

import { HttpClientModule, HttpClient, HttpHeaders } from '@angular/common/http';
import { DialogFuzzLive } from './dialog-fuzz-live.component';

import {
    MatDialog,
    MatDialogRef,
    MatSnackBar,
    MAT_DIALOG_DATA
} from '@angular/material';

export interface FuzzData {
    name: string;
    detail: boolean;
    show: boolean;
    output_path: string;
    state_info: string;
    fuzz_ip: string;
    start_time: string;
    last_update: string;
    execs_per_sec: string;
    unique_crashes: string;
    unique_hangs: string;
    loading: boolean;
}

export interface UILiveData {
    name: string;
    ip: string;
}

@Component({
    selector: 'app-fuzz',
    templateUrl: './fuzz.component.html',
    styleUrls: ['./fuzz.component.css'],
})


export class FuzzComponent implements OnInit {
    consoledata: Array<FuzzData> = [];
    guidatas: Array<UILiveData> = [];
    constructor(
        private http: HttpClient,
        public dialog: MatDialog) {
    }

    ngOnInit(): void {
        // let url_s: string = "http://10.103.129.79/test/state/heartbeat";
        let url_fxcore_s: string = "http://10.103.129.9:9090/testserver?label=linux_fuzz";
        this.http
            .get(url_fxcore_s).subscribe(resp => {
                let data = resp['data'];
                if (data == null) {
                    return;
                }
                // console.info(resp['data']);
                let current_machine: any = {};
                for (let i = 0; i < data.length; i++) {
                    // console.info(data[i]);
                    let ip: string = data[i].ip;
                    this.getData(ip);
                }
            });
            // let uidata: UILiveData = {
            //     ip: '10.203.18.113',
            //     name: 'tem'
            // }
            // this.guidatas.push(uidata);
        let url_gui_s: string = "test/state/heartbeat";
        this.http
            .get(url_gui_s).subscribe(resp => {
                let data = resp['data'];
                if (data == null) {
                    return;
                }
                // console.info(resp['data']);
                let current_machine: any = {};
                for (let i = 0; i < data.length; i++) {
                    // console.info(data[i]);
                    let ip: string = data[i].machine.split('_')[1];
                    this.getData(ip + ':9091');
                }
            });

    }

    getData(fuzz_ip: string): void {
        let fuzz_url_s: string = 'http://' + fuzz_ip + '/fuzz-data';
        this.http
            .get(fuzz_url_s).subscribe(resp => {
                let data = resp['data'];
                for (let i = 0; i < data.length; i++) {
                    // data[i] = 'UI_LIVE_LOCALTEST';
                    console.info(data);
                    let ui_i: number = data[i].indexOf('UI_LIVE_');
                    if (ui_i != -1) {
                        let uidata: UILiveData = {
                            ip: fuzz_ip.substring(0, fuzz_ip.length - 5),
                            name: data[i].substring(ui_i + 8)
                        }
                        this.guidatas.push(uidata);
                    } else {
                        let fuzzdata: FuzzData = {
                            name: data[i],
                            detail: false,
                            show: true,
                            loading: false,
                            output_path: '',
                            fuzz_ip: fuzz_ip,
                            start_time: '',
                            last_update: '',
                            execs_per_sec: '',
                            unique_crashes: '',
                            unique_hangs: '',
                            state_info: ''
                        };
                        this.consoledata.push(fuzzdata);
                        fuzzdata.name = data[i];
                    }

                }


            });
    }

    closeCard(temthis: any): void {
        console.info(temthis);
        temthis.show = false;
        // console.info()
        //this.showCard = false;
    }

    showDetail(temthis: any): void {
        temthis.loading = true;
        temthis.detail = true;
        let url_bmp_s: string = "http://" + temthis.fuzz_ip + "/fuzz-data?type=bitmap&name=" + temthis.name;
        this.http
            .get(url_bmp_s).subscribe(bmpresp => {
                let url_state_s: string = "http://" + temthis.fuzz_ip + "/fuzz-data?type=stats&name=" + temthis.name;
                this.http
                    .get(url_state_s).subscribe(stateresp => {
                        temthis.state_info = stateresp['data'];
                    });

                let url_plot_s: string = "http://" + temthis.fuzz_ip + "/fuzz-data?type=plot&name=" + temthis.name;
                this.http
                    .get(url_plot_s).subscribe(plotresp => {
                        let url_s: string = "test/fuzz/chart-data";
                        const body = 'name=' + temthis.name + '&bitmap=' + encodeURI(bmpresp['data']).replace(/\+/g, '%2B') + '&plot=' + plotresp['data'];
                        // console.info(body);
                        this.http
                            .post(url_s, body, {
                                headers: new HttpHeaders().set('Content-Type', 'application/x-www-form-urlencoded'),
                            }).subscribe(resp => {
                                temthis.output_path = resp['data'];
                                temthis.loading = false;
                                // this.dialogRef.close();
                            })
                    });

            });
        let url_stats_s: string = "http://" + temthis.fuzz_ip + "/fuzz-data?type=stats&name=" + temthis.name;
        console.info(url_stats_s);
        // url_stats_s = "http://"+'10.103.128.216:9091'+"/fuzz-data?type=stats&name=" + 'ConvertToPDF_PNG';
        this.http
            .get(url_stats_s).subscribe(statsresp => {
                console.info(statsresp);
                if (statsresp['ret'] == -2) {
                    temthis.last_update = "Interrupted";
                    return;
                }
                let stats: string = this.b64DecodeUnicode(statsresp['data'])
                let stats_sp = stats.split('\n');
                for (let i = 0; i < stats_sp.length; i++) {
                    // console.info(stats_sp[i]);
                    let stat = stats_sp[i];
                    if (stat.indexOf('start_time') != -1) {
                        temthis.start_time = this.timeConverter(+stat.split(':')[1]);
                    } else if (stat.indexOf('last_update') != -1) {
                        temthis.last_update = this.timeConverter(+stat.split(':')[1]);
                    } else if (stat.indexOf('execs_per_sec') != -1) {
                        temthis.execs_per_sec = stat.split(':')[1];
                    } else if (stat.indexOf('unique_crashes') != -1) {
                        temthis.unique_crashes = stat.split(':')[1];
                    } else if (stat.indexOf('unique_hangs') != -1) {
                        temthis.unique_hangs = stat.split(':')[1];
                    }
                }
            });
    }

    timeConverter(UNIX_timestamp: number): string {
        var a = new Date(UNIX_timestamp * 1000);
        // return a.toLocaleTimeString();
        var months = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
        var year = a.getFullYear();
        var month = months[a.getMonth()];
        var date = a.getDate();
        var hour = a.getHours() < 10 ? '0' + a.getHours() : a.getHours();
        var min = a.getMinutes() < 10 ? '0' + a.getMinutes() : a.getMinutes();
        var sec = a.getSeconds() < 10 ? '0' + a.getSeconds() : a.getSeconds();

        var time = date + ' ' + month + ' ' + hour + ':' + min + ':' + sec;
        return time;
    }

    b64DecodeUnicode(str: string): string {
        // Going backwards: from bytestream, to percent-encoding, to original string.
        return decodeURIComponent(atob(str).split('').map(function (c) {
            return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
        }).join(''));
    }

    // openLive(ip: string, name: string): void {
    //     console.info("OPEN");
    //     this.openLiveDialog(ip);
    // }

    openLiveDialog(ip: string, name: string): void {
        // console.info(msg);
        let dialogRef = this.dialog.open(DialogFuzzLive, {
            //   width: '250px',
            height: '640px',
            width: '800px',
            data: { ip: ip, name: name }
        });

        dialogRef.afterClosed().subscribe(result => {
        });
    }

}

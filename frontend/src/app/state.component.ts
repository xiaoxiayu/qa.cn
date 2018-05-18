import { Component, OnInit, OnDestroy, ViewChild, ChangeDetectorRef, Inject } from '@angular/core';
import {
    MatPaginator,
    MatTableDataSource,
    MatDialog,
    MatDialogRef,
    MatSnackBar,
    MAT_DIALOG_DATA,
    MatSort,
    MatChipInputEvent
} from '@angular/material';
import { ENTER, COMMA } from '@angular/cdk/keycodes';
import { HttpClientModule, HttpClient, HttpHeaders } from '@angular/common/http';
import { $WebSocket, WebSocketSendMode } from 'angular2-websocket/angular2-websocket';
import 'rxjs/add/operator/switchMap';
import { ActivatedRoute, ParamMap } from '@angular/router';

import { DialogEditMachine } from './dialog-edit-machine.component';

import { Observable, Subscriber } from 'rxjs/Rx';

import { StateService } from './state.service'


@Component({
    selector: 'app-state',
    templateUrl: './state.component.html',
    styleUrls: ['./state.component.css']
})

export class StateComponent implements OnInit, OnDestroy {
    panelOpenState: boolean = true;
    selectedIndex: number = 0;
    fxcore_running_data: RunningState[] = [];
    fxcore_machine_data: MachineState[] = [];
    set_running_data: RunningState[] = [];
    set_machine_data: MachineState[] = [];
    set_machine_ui_autolinked: any[] = [];
    set_testsummary_data: TestSummary[] = [];
    fxcore_state_loading: boolean = true;
    fxcore_machine_loading: boolean = true;
    set_interval_instance: any = null;
    fxcore_running_interval_instance: any = null;
    fxcore_machine_interval_instance: any = null;


    // For Search And Sort.
    set_machine_map: any = {};
    fxcore_machine_map: any = {};

    displayedColumns = ['position', 'name', 'creator', 'progress', 'status', 'tools'];
    fxcore_machine_displayedColumns = ['position', 'name', 'test', 'info', 'status', 'message'];
    set_machine_displayedColumns = ['position', 'name', 'test', 'info', 'status', 'ui', 'message'];
    set_running_displayedColumns = ['position', 'name', 'creator', 'progress', 'status', 'tools'];
    set_testsummary_displayedColumns = ['position', 'name', 'testindex', 'ci', 'status', 'result'];

    dataSource = new MatTableDataSource<RunningState>(this.fxcore_running_data);
    set_machine_dataSource = new MatTableDataSource<MachineState>(this.set_machine_data);
    set_running_dataSource = new MatTableDataSource<RunningState>(this.set_running_data);
    set_testsummary_dataSource = new MatTableDataSource<TestSummary>(this.set_testsummary_data);
    fxcore_machine_dataSource = new MatTableDataSource<MachineState>(this.fxcore_machine_data);

    // @ViewChild(MatPaginator) paginator: MatPaginator;
    // @ViewChild(MatPaginator) set_machine_paginator: MatPaginator;
    @ViewChild('paginator') paginator: MatPaginator;
    @ViewChild('fxcore_machine_paginator') fxcore_machine_paginator: MatPaginator;
    @ViewChild('set_machine_paginator') set_machine_paginator: MatPaginator;
    @ViewChild('set_running_paginator') set_test_paginator: MatPaginator;
    @ViewChild('set_testsummary_paginator') set_testsummary_paginator: MatPaginator;

    visible: boolean = true;
    selectable: boolean = true;
    removable: boolean = true;
    addOnBlur: boolean = true;

    // Enter, comma
    separatorKeysCodes = [ENTER, COMMA];

    constructor(
        private http: HttpClient,
        public dialog: MatDialog,
        private route: ActivatedRoute,
        private service: StateService,
        public snackBar: MatSnackBar,
    ) { };
    /**
     * Set the paginator after the view init since this component will
     * be able to query its view for the initialized paginator.
     */
    ngAfterViewInit() {
        this.dataSource.paginator = this.paginator;
        this.fxcore_machine_dataSource.paginator = this.fxcore_machine_paginator;
        this.set_machine_dataSource.paginator = this.set_machine_paginator;
        this.set_running_dataSource.paginator = this.set_test_paginator;
        this.set_testsummary_dataSource.paginator = this.set_testsummary_paginator;
        //this.fxcore_machine_dataSource.sort = this.sort;
    }

    ngOnDestroy(): void {
        if (this.fxcore_running_interval_instance != null) {
            clearInterval(this.fxcore_running_interval_instance);
        }
        if (this.fxcore_machine_interval_instance != null) {
            clearInterval(this.fxcore_machine_interval_instance);
        }
        if (this.set_interval_instance != null) {
            clearInterval(this.set_interval_instance);
        }
    }

    ngOnInit(): void {
        let id = +this.route.snapshot.paramMap.get('id');
        if (id == 0) {
            this.fxcore_machine();
            this.fxcore_running();
            this.fxcore_running_interval_instance = setInterval(() => { this.fxcore_running(); }, 5000);
            this.dataSource.data = this.fxcore_running_data;
        } else if (id == 1) {
            this.gui_update();
            this.set_interval_instance = setInterval(() => { this.gui_update(); }, 5000);
            this.gui_update_summary();
            this.set_machine_dataSource.data = this.set_machine_data;
        }

        this.route.paramMap
            .switchMap((params: ParamMap) =>
                this.service.updateTab(params.get('id')))
            .subscribe(select_id => this.selectedIndex = +select_id);
    }

    add_fxcore_label(position: number, ip: string, event: MatChipInputEvent): void {
        let input = event.input;
        let value = event.value;

        // Add our fruit
        if ((value || '').trim()) {
            //this.fruits.push({ name: value.trim() });
            let labels = this.fxcore_machine_data[position].labels;
            labels.push(value);

            let url_s = "http://" + ip + '/label'
            const body = "set=" + labels.toString();
            this.http
                .put(url_s, body, {
                    headers: new HttpHeaders().set('Content-Type', 'application/x-www-form-urlencoded'),
                }).subscribe(resp => {
                    console.info(resp);
                    if (resp['ret'] != 0) {
                        setTimeout(() => this.snackBar.open(('Add "' + value + '" for ' + ip + ' failed.')))
                        return;
                    }
                    this.fxcore_machine_data[position].labels = labels;
                    this.fxcore_machine_dataSource.data = this.fxcore_machine_data;
                })
        }

        // Reset the input value
        if (input) {
            input.value = '';
        }
    }

    remove_fxcore_label(position: number, ip: string, label: string): void {
        let labels = this.fxcore_machine_data[position].labels;
        let index = this.fxcore_machine_data[position].labels.indexOf(label);

        if (index >= 0) {
            labels.splice(index, 1);
            let url_s = "http://" + ip + '/label'
            const body = "set=" + labels.toString();
            this.http
                .put(url_s, body, {
                    headers: new HttpHeaders().set('Content-Type', 'application/x-www-form-urlencoded'),
                }).subscribe(resp => {
                    console.info(resp);
                    if (resp['ret'] != 0) {
                        setTimeout(() => this.snackBar.open(('Remove "' + label + '" for ' + ip + ' failed.')))
                        return;
                    }
                    this.fxcore_machine_data[position].labels = labels;
                    this.fxcore_machine_dataSource.data = this.fxcore_machine_data;
                })
        }
        this.fxcore_machine_dataSource.data = this.fxcore_machine_data;
    }

    tabChanged(event: any) {
        if (this.fxcore_running_interval_instance != null) {
            clearInterval(this.fxcore_running_interval_instance);
        }
        if (this.fxcore_machine_interval_instance != null) {
            clearInterval(this.fxcore_machine_interval_instance);
        }
        if (this.set_interval_instance != null) {
            clearInterval(this.set_interval_instance);
        }
        if (event.tab.textLabel == 'Fxcore') {
            this.fxcore_machine();
            this.fxcore_running();
            this.dataSource.data = this.fxcore_running_data;
        } else if (event.tab.textLabel == 'SET') {
            this.gui_update();
            this.set_interval_instance = setInterval(() => { this.gui_update(); }, 5000);
            this.gui_update_summary();
            this.set_machine_dataSource.data = this.set_machine_data;
        }

    }

    applyFilter(filterValue: string) {
        filterValue = filterValue.trim();
        filterValue = filterValue.toLowerCase();
        this.dataSource.filter = filterValue;
    }

    applyFxcoreMachineFilter(filterValue: string) {
        filterValue = filterValue.trim();
        filterValue = filterValue.toLowerCase();
        this.fxcore_machine_dataSource.filter = filterValue;
    }

    applySETMachineFilter(filterValue: string) {
        filterValue = filterValue.trim();
        filterValue = filterValue.toLowerCase();
        this.set_machine_dataSource.filter = filterValue;
    }

    applySETRunningFilter(filterValue: string) {
        filterValue = filterValue.trim();
        filterValue = filterValue.toLowerCase();
        this.set_running_dataSource.filter = filterValue;
    }

    applySETTestSummaryFilter(filterValue: string) {
        filterValue = filterValue.trim();
        filterValue = filterValue.toLowerCase();
        this.set_testsummary_dataSource.filter = filterValue;
    }

    machine_message_update(ip: string, name: string): void {
        //let url_s: string = "http://10.103.129.79/test/state/gui-machine-msg?ip=" + ip;
        let url_s: string = "test/state/gui-machine-msg?ip=" + ip;
        this.http
            .get(url_s).subscribe(resp => {
                for (let i = 0; i < this.set_machine_data.length; i++) {
                    if (this.set_machine_data[i].name == ip) {
                        this.set_machine_data[i]['message'] = resp['message'];
                        this.set_machine_data[i]['message_update_time'] = resp['updatetime'];
                        this.set_machine_data[i]['purpose'] = resp['purpose'];
                        this.set_machine_dataSource.data = this.set_machine_data;
                        break;
                    }
                }

            });
    }

    _compare_unixtime(a: any, b: any): number {
        let _a: number = (new Date(a.updatetime)).getTime() / 1000;
        let _b: number = (new Date(b.updatetime)).getTime() / 1000;
        if (_a > _b)
            return -1;
        if (_a < _b)
            return 1;
        return 0;
    }

    gui_update_summary(): void {
        // let url_s: string = "http://10.103.129.79/test/state/guisummary";
        let url_s: string = "test/state/guisummary";
        this.http
            .get(url_s).subscribe(resp => {
                console.info(resp);
                let data = resp['data'];
                if (data == null) {
                    return;
                }

                for (let i = 0; i < data.length; i++) {
                    let result_s: string = data[i].Result;
                    if (data[i].Resulturl != "") {
                        result_s = '<a href="' + data[i].Resulturl + '">' + data[i].Result + '</a>';
                    }

                    // (new Date(data[i].Updatetime).getTime()/1000)

                    this.set_testsummary_data.push({
                        position: i,
                        name: data[i].Testname,
                        testindex: data[i].Testindex,
                        ci: data[i].Ci,
                        status: data[i].Status,
                        result: result_s,
                        updatetime: data[i].Updatetime,
                        resulturl: data[i].Resulturl
                    });
                }
                this.set_testsummary_data.sort(this._compare_unixtime);

                for (let i = 0; i < this.set_testsummary_data.length; i++) {
                    this.set_testsummary_data[i].position = i;
                }

                this.set_testsummary_dataSource.data = this.set_testsummary_data;
            })

    }

    _check_set_running_exists(testname: string): boolean {
        for (let i = 0; i < this.set_running_data.length; i++) {
            let data = this.set_running_data[i];
            if (data.name == testname) {
                return true;
            }
        }
        return false;
    }

    _get_set_machine_index(ip: string): number {
        for (let i = 0; i < this.set_machine_data.length; i++) {
            let data = this.set_machine_data[i];
            if (data.name == ip) {
                return i;
            }
        }
        return -1;
    }

    _get_set_running_index(name: string): number {
        for (let i = 0; i < this.set_running_data.length; i++) {
            let data = this.set_running_data[i];
            if (data.name == name) {
                return i;
            }
        }
        return -1;
    }

    changeUI(element: MachineState, status: string, resolution: string): void {
        if (status == 'unlink') {
            let url_s: string = 'http://10.203.13.116:9091/rdplink/' + element.name;
            let options: any = {};
            this.http.delete(url_s, options).subscribe(resp => {
                if (resp['ret'] == 0) {
                    setTimeout(() => this.snackBar.open(('unlink: ' + name + ' ok.'), '', { duration: 5000 }));
                    element.uilink_status = status;
                    return;
                }
                setTimeout(() => this.snackBar.open(('unlink: ' + name + ' failed.'), '', { duration: 5000 }));
            });
        } else if (status == 'autolink') {
            let url_s: string = 'http://10.203.13.116:9091/rdplink';
            let body: string = 'ip=' + element.name;
            let screen_size: Array<string> = resolution.split('x');
            if (screen_size.length > 1) {
                body = 'ip=' + element.name + '&width=' + screen_size[0] + '&height=' + screen_size[1];
            }
            this.http
                .post(url_s, body, {
                    headers: new HttpHeaders().set('Content-Type', 'application/x-www-form-urlencoded'),
                }).subscribe(resp => {
                    if (resp['ret'] == 0) {
                        setTimeout(() => this.snackBar.open(('autolink: ' + name + ' ok.'), '', { duration: 5000 }));
                        element.uilink_status = status;
                        // element.uiinfo = 'UI: 10.203.13.116';
                        element.screen_resolution = resolution;
                        return;
                    }
                    setTimeout(() => this.snackBar.open(('autolink: ' + name + ' failed.'), '', { duration: 5000 }));
                })
        }
    }

    _check_ip_ui_autolink(ip: string): string {
        for (let i = 0; i < this.set_machine_ui_autolinked.length; i++) {
            let data = this.set_machine_ui_autolinked[i];
            //console.info('***LinkedUI:'+data);
            if (data['ip'] == ip) {
                return data['res'];
            }
        }
        return '';
    }

    gui_update(): void {
        // this.set_machine_data.push({
        //     position: 0,
        //     name: 'ip',
        //     info: '',
        //     update_time: 'this.set_machine_map[name].time',
        //     test: 'this.set_machine_map[name].testname',
        //     // test: this.set_machine_map[name].testinfo,
        //     status: 'this.set_machine_map[name].status',
        //     message: 'this.set_machine_map[name].message',
        //     message_update_time: 'this.set_machine_map[name].message_update_time',
        //     purpose: 'this.set_machine_map[name].purpose',
        //     labels: [],
        //     envinfos: [],
        //     uiinfo: '',
        //     uilink_status: 'unlink',
        //     screen_resolution: '',
        // });

        //let url_s: string = "http://10.103.129.79/test/state/heartbeat";
        this._machine_UI_info('set');
        let url_s: string = "test/state/heartbeat";
        // console.info('herer:', url_s);

        this.http
            .get(url_s).subscribe(resp => {
                let data = resp['data'];
                if (data == null) {
                    return;
                }

                let current_machine: any = {};
                for (let i = 0; i < data.length; i++) {
                    // console.info(data[i]);
                    this.set_machine_map[data[i].machine] = data[i];
                    current_machine[data[i].machine] = '';
                }

                // Get offline machine.
                let offline_machine: Array<string> = [];
                for (var machinename in this.set_machine_map) {
                    if (machinename in current_machine) {
                        continue;
                    } else {
                        offline_machine.push(machinename);
                    }
                }

                // Remove offline machine.
                for (let i = 0; i < offline_machine.length; i++) {
                    delete this.set_machine_map[offline_machine[i]];
                }

                // Set table data.
                //this.set_machine_data = [];
                var keys = Object.keys(this.set_machine_map);
                keys.sort();


                let current_running: any = {};
                for (let i = 0; i < keys.length; i++) {
                    let name = keys[i];
                    let ip: string = name.split('_')[1];

                    if (this.set_machine_map[name].testname != "") {
                        current_running[this.set_machine_map[name].testname] = '';

                        let color: string = 'primary';
                        let progress_s: string = this.set_machine_map[name].testinfo.split(':')[0];
                        let progress_tem = progress_s.split('/');
                        let progress_i: number = +((+(progress_tem[0]) / +(progress_tem[1])) * 100).toFixed(0);
                        // console.info('Progress:', p)
                        let state_s: string = 'Initialization';
                        if (progress_i >= 1) {
                            state_s = 'Running';
                        }

                        if (!this._check_set_running_exists(this.set_machine_map[name].testname)) {
                            this.set_running_data.push({
                                position: 0,
                                createtime: this.set_machine_map[name].create_time,
                                name: this.set_machine_map[name].testname,
                                creator: this.set_machine_map[name].creator,
                                progress: progress_i,
                                status: state_s,
                                color: color,
                                create_url: '',
                                status_info: this.set_machine_map[name].testinfo,
                                tools: '',
                                executor: '',
                                executor_pid: ''
                            });
                        } else {
                            let test_i: number = this._get_set_running_index(this.set_machine_map[name].testname);
                            this.set_running_data[test_i].progress = progress_i;
                            this.set_running_data[test_i].status = state_s;
                            this.set_running_data[test_i].status_info = this.set_machine_map[name].testinfo;
                        }
                    }
                    let machine_i: number = this._get_set_machine_index(ip);
                    if (machine_i == -1) {
                        this.set_machine_data.push({
                            position: i,
                            name: ip,
                            info: '',
                            update_time: this.set_machine_map[name].time,
                            test: this.set_machine_map[name].testname,
                            // test: this.set_machine_map[name].testinfo,
                            status: this.set_machine_map[name].status,
                            message: this.set_machine_map[name].message,
                            message_update_time: this.set_machine_map[name].message_update_time,
                            purpose: this.set_machine_map[name].purpose,
                            labels: [],
                            envinfos: [],
                            uiinfo: '',
                            uilink_status: 'unlink',
                            screen_resolution: '',
                        });

                    } else {
                        this.set_machine_data[machine_i].update_time = this.set_machine_map[name].time;
                        this.set_machine_data[machine_i].test = this.set_machine_map[name].testname;
                        this.set_machine_data[machine_i].status = this.set_machine_map[name].status;
                        // this.set_machine_data[machine_i].message = this.set_machine_map[name].message;
                        this.set_machine_data[machine_i].message_update_time = this.set_machine_map[name].message_update_time;
                        this.set_machine_data[machine_i].purpose = this.set_machine_map[name].purpose;
                    }

                    // Machine message update.
                    this.machine_message_update(ip, name);
                    let update_ui_state: boolean = true;
                    let screen_res: string = this._check_ip_ui_autolink(ip)
                    if (screen_res != '') {
                        this.set_machine_data[machine_i].uiinfo = 'AutoUI: 10.203.13.116';
                        this.set_machine_data[machine_i].screen_resolution = screen_res;
                        update_ui_state = false;
                    }
                    this._machine_update_info('set', machine_i, ip + ':9091', update_ui_state);

                }

                // Get stopped running
                let stoped_running: Array<number> = [];
                for (var i = 0; i < this.set_running_data.length; i++) {
                    let running_data: RunningState = this.set_running_data[i];
                    if (running_data.name in current_running) {
                        continue;
                    } else {
                        stoped_running.push(i);
                    }
                }

                // Remove stopped running.
                for (let i = 0; i < stoped_running.length; i++) {
                    //delete this.set_running_data[offline_machine[i]];
                    let stopped_i: number = stoped_running[i];
                    this.set_running_data.splice(stopped_i, 1);
                }

                // Update Test Table.
                this.set_machine_dataSource.data = this.set_machine_data;

                this.set_running_dataSource.data = this.set_running_data;


            },
                resp => {
                    console.info('ERROR:', resp);
                });
    }

    stopSetRunning(creator_url: string, name: string): void {
        let stop_ci_url: string = 'cistop';
        let stop_body: string = 'ciurl=' + creator_url;
        this.http
            .post(stop_ci_url, stop_body, {
                headers: new HttpHeaders().set('Content-Type', 'application/x-www-form-urlencoded'),
            }).subscribe(resp => {
                setTimeout(() => {
                    let url_s: string = 'http://10.103.129.147:9090/end';
                    let body: string = 'action=force&test=' + name;
                    this.http
                        .post(url_s, body, {
                            headers: new HttpHeaders().set('Content-Type', 'application/x-www-form-urlencoded'),
                        }).subscribe(resp => {
                            console.info(resp);
                            setTimeout(() => this.snackBar.open(('Task: ' + name + ' stoped.'), '', { duration: 5000 }));
                        })
                }, 3000)


                setTimeout(() => this.snackBar.open((creator_url + ' stoped. Wait test process stop...'), '', { duration: 10000 }));
            })


    }

    _machine_UI_info(group: string): void {
        let rdpurl_s: string = "http://10.203.13.116:9091/rdplink-connected";
        this.http
            .get(rdpurl_s).subscribe(resp => {
                if (resp["ret"] == 0) {
                    // let data = resp["data"];
                    this.set_machine_ui_autolinked = resp["data"];
                    // for (let i: number = 0; i < data.length; i++) {
                        // this.set_machine_ui_autolinked = data[i]['ip'];
                        
                    // }
                    

                }
            },
                err_resp => {
                    // this.set_machine_data[position].envinfos = envinfos;
                });
    }

    _machine_update_info(group: string, position: number, ip: string, update_ui_status: boolean): void {
        if (group == 'set') {
            if (update_ui_status) {
                let rdpurl_s: string = 'http://' + ip + '/rdplink'
                this.http
                    .get(rdpurl_s).subscribe(resp => {
                        if (resp["ret"] == 0) {
                            // envinfos.push('UI: ' + resp["ip"]);
                            if (resp["ip"] == '') {
                                this.set_machine_data[position].uilink_status = 'unlink';
                            } else {
                                this.set_machine_data[position].uilink_status = 'userlink';
                            }
                            this.set_machine_data[position].uiinfo = 'UI: ' + resp["ip"];
                        }
                    },
                        err_resp => {
                            this.set_machine_data[position].uilink_status = 'unlink';
                            this.set_machine_data[position].uiinfo = 'UI: Get Info Failed.';
                        });
            } else {
                this.set_machine_data[position].uilink_status = 'autolink';
            }
        }

        let url_s: string = 'http://' + ip + '/hardware?type=memory'
        this.http
            .get(url_s).subscribe(resp => {
                // console.info(resp);
                let data = resp['data'];
                let envinfos: Array<string> = [];
                if (data.os == 'windows') {
                    envinfos.push('OS: ' + data.os + ' ' + data.core);
                } else {
                    envinfos.push('OS: ' + data.os + ' ' + data.platform);
                }
                envinfos.push('Hostname: ' + data.hostname);

                let info_detail = 'Cpus: ' + data.cpus +
                    ';\n Memorys(MB): Total: ' + data.memory.total / 1048576 +
                    ' UsedPercent: ' + data.memory.usedPercent.toFixed(0) + '%';

                if (group == 'fxcore') {
                    this.fxcore_machine_data[position].info = info_detail;
                    this.fxcore_machine_data[position].envinfos = envinfos;
                } else if (group == 'set') {
                    this.set_machine_data[position].envinfos = envinfos;
                    this.set_machine_dataSource.data = this.set_machine_data;
                }

                //this.fxcore_machine_dataSource.data = this.fxcore_machine_data;

            },
                resp => {
                    console.info('Can not link:' + resp.url);
                }
            );

    }

    _update_fxcore_machine_array_data(data: any): void {
        for (let i = 0; i < this.fxcore_machine_data.length; i++) {
            if (data.ip == this.fxcore_machine_data[i].name) {
                this.fxcore_machine_data[i].status = data.taskcount;
                this.fxcore_machine_data[i].labels = data.Label.split(',');
                return;
            }
        }
        this.fxcore_machine_data.push({
            position: this.fxcore_machine_data.length,
            name: data.ip,
            info: '',
            update_time: '',
            test: '',
            status: data.taskcount,
            message: '',
            message_update_time: '',
            purpose: data.Label,
            labels: data.Label.split(','),
            envinfos: [],
            uiinfo: '',
            uilink_status: 'unlink',
            screen_resolution: '',
        });
    }

    _do_fxcore_machine_update(): void {
        this.fxcore_machine_dataSource.data = this.fxcore_machine_data;
        let url_s: string = 'http://10.103.129.9:9090/testserver'; // Only Mac TestServer Now.
        this.http
            .get(url_s).subscribe(resp => {
                // console.info(resp);
                let data = resp['data'];
                console.info(data);
                //let datamap = {};
                let current_machine = {};
                for (let i = 0; i < data.length; i++) {
                    current_machine[data[i].ip] = '';
                    this.fxcore_machine_map[data[i].ip] = data[i];
                }

                // Get offline machine.
                let offline_machine: Array<string> = [];
                for (var machinename in this.fxcore_machine_map) {
                    if (machinename in current_machine) {
                        continue;
                    } else {
                        offline_machine.push(machinename);
                    }
                }

                // Remove offline
                for (let i = 0; i < offline_machine.length; i++) {
                    delete this.fxcore_machine_map[offline_machine[i]];
                }

                var keys = Object.keys(this.fxcore_machine_map);
                keys.sort();

                for (let i = 0; i < keys.length; i++) {
                    let name = keys[i];

                    let _data = this.fxcore_machine_map[name];
                    this._update_fxcore_machine_array_data(_data);
                    this._machine_update_info('fxcore', i, name, false);
                }
                this.fxcore_machine_loading = false;
                this.fxcore_machine_dataSource.data = this.fxcore_machine_data;

            })
    }

    fxcore_machine(): void {
        this.fxcore_machine_data = [];
        this._do_fxcore_machine_update();
        this.fxcore_machine_interval_instance = setInterval(() => {
            this._do_fxcore_machine_update();

        }, 10000);

    }

    _check_fxcorerunning_data_exists(name: string): number {
        for (let j = 0; j < this.fxcore_running_data.length; j++) {
            if (this.fxcore_running_data[j].name == name) {
                return j;
            }
        }
        return -1;
    }

    fxcore_running(): void {
        let url_s: string = "http://10.103.129.80:32457/set?key=FXQA_TESTING"
        this.http
            .get(url_s).subscribe(resp => {
                let tests = resp['val'];
                if (tests == null) {
                    return;
                }

                for (let i = 0; i < tests.length; i++) {
                    // alert(tests[i]);

                    let test_url: string = "http://10.103.129.80:32457/hash?key=FXQA-" + tests[i];
                    this.http
                        .get(test_url).subscribe(resp => {
                            let color: string = 'primary';
                            if (resp['status'] == 'Success' || resp['status'] == 'Initialization') {
                                color = 'primary';
                            } else if (resp['status'] == 'Unknown') {
                                color = 'warn';
                            } else {
                                color = 'accent';
                            }

                            let exists_i  = this._check_fxcorerunning_data_exists(tests[i]);
                            if (exists_i == -1) {
                                this.fxcore_running_data.push({
                                    position: i,
                                    createtime: resp['starttime'],
                                    name: tests[i],
                                    creator: resp['creator'],
                                    progress: resp['progress'],
                                    status: resp['status'],
                                    color: color,
                                    create_url: resp['ci'],
                                    status_info: resp['info'],
                                    tools: '',
                                    executor: resp['executor'],
                                    executor_pid: resp['executor_pid']
                                });
                            } else {
                                this.fxcore_running_data[exists_i].color = color;
                                this.fxcore_running_data[exists_i].progress = resp['progress'];
                                this.fxcore_running_data[exists_i].status = resp['status'];
                                this.fxcore_running_data[exists_i].executor = resp['executor'];
                                this.fxcore_running_data[exists_i].executor_pid = resp['executor_pid'];
                            }
                            
                            this.dataSource.data = this.fxcore_running_data;
                        })
                }
                this.fxcore_state_loading = false;

            },
                resp => {
                    console.info('ERROR:', resp);
                });
        // let ws_url = document.location.href.replace('http://', 'ws://');
        // ws_url = ws_url.slice(0, ws_url.indexOf('/test/state')) + '/test/state/_data';

        // let ws = new $WebSocket(ws_url);
        // let position_i: number = 0;

        // // set received message stream
        // ws.getDataStream().subscribe(
        //     (msg) => {
        //         // console.log(msg);
        //         this.fxcore_state_loading = false;
        //         var data = JSON.parse(msg.data);
        //         console.info(data.type);
        //         if (data.type == 'createall_end') {
        //             this.dataSource.data = this.fxcore_running_data;

        //             this.fxcore_running_interval_instance = setInterval(() => {
        //                 //ws.send('updatealldata');
        //                 ws.send("updatealldata").subscribe(
        //                     (msg) => {
        //                         console.log("next", msg.data);

        //                     },
        //                     (msg) => {

        //                         console.log("error4", msg);
        //                     },
        //                     (msg) => {
        //                         console.info(msg);
        //                         // this.snackBar.open('Update failed. Please try refreshing the page.');
        //                         //setTimeout(() => this.snackBar.open(('Update failed. Please try refreshing the page.')))
        //                     }
        //                 );
        //                 // this.fxcore_running_data.push({position: 1, name: 'Hydrogen', creator: '', progress: 1.0079, tools: 'H'});
        //                 // this.dataSource.data = this.fxcore_running_data;
        //                 // console.info(this.dataSource.data);
        //             }, 20000);
        //         } else if (data.type == 'create') {
        //             // console.info(data);
        //             if (data.start_time == '') {
        //                 data.start_time = 'Start time not set.';
        //             }
        //             let color: string = 'primary';
        //             if (data.status == 'Success' || data.status == 'Initialization') {
        //                 color = 'primary';
        //             } else if (data.status == 'Unknown') {
        //                 color = 'warn';
        //             } else {
        //                 color = 'accent';
        //             }

        //             this.fxcore_running_data.push({
        //                 position: position_i,
        //                 createtime: data.start_time,
        //                 name: data.task,
        //                 creator: data.creator,
        //                 progress: data.progress,
        //                 status: data.status,
        //                 color: color,
        //                 create_url: data.create_url,
        //                 status_info: data.status_info,
        //                 tools: '',
        //                 executor: data.executor,
        //                 executor_pid: data.executor_pid
        //             });
        //             // console.info(this.fxcore_running_data);
        //             position_i++;
        //         } else if (data.type == 'updatealldata_start') {
        //             // console.info('CurrentCount:', data.count, 'OriginCount:', this.fxcore_running_data.length);
        //             this.fxcore_state_loading = true;
        //             if (data.count == this.fxcore_running_data.length) {
        //                 ws.send('doupdate').subscribe(
        //                     (msg) => {
        //                         console.log("next", msg.data);

        //                     },
        //                     (msg) => {

        //                         console.log("error3", msg);
        //                     },
        //                     () => {
        //                         console.log("complete");
        //                     }
        //                 );
        //             } else {
        //                 this.fxcore_running_data = [];
        //             }
        //         } else if (data.type == 'update') {
        //             this.fxcore_state_loading = true;
        //             let color: string = 'primary';
        //             if (data.status == 'Success' || data.status == 'Initialization') {
        //                 color = 'primary';
        //             } else if (data.status == 'Unknown') {
        //                 color = 'warn';
        //             } else {
        //                 color = 'accent';
        //             }
        //             this.fxcore_running_data.push({
        //                 position: position_i,
        //                 createtime: data.start_time,
        //                 name: data.task,
        //                 creator: data.creator,
        //                 progress: data.progress,
        //                 status: data.status,
        //                 color: color,
        //                 create_url: data.create_url,
        //                 status_info: data.status_info,
        //                 tools: '',
        //                 executor: data.executor,
        //                 executor_pid: data.excutor_pid
        //             });

        //         } else if (data.type == 'updateall_end') {
        //             this.fxcore_state_loading = false;
        //             this.dataSource.data = this.fxcore_running_data;
        //         }

        //     },
        //     (msg) => {
        //         console.log("error2", msg);
        //     },
        //     () => {
        //         console.log("complete");
        //     }
        // );

        // //     // send with default send mode (now default send mode is Observer)
        // ws.send("createalldata").subscribe(
        //     (msg) => {
        //         console.log("next", msg.data);

        //     },
        //     (msg) => {
        //         console.info(msg);
        //         if (msg == 'Socket connection has been closed') {
        //             //setTimeout(() => this.snackBar.open('Running: Connet failed. Please try refreshing the page.'));
        //         }

        //     },
        //     () => {
        //         setTimeout(() => this.snackBar.open('Get running data failed. Please try refreshing the page.', '', { duration: 3000 }))
        //     }
        // );

    }

    openDialog(ip: string, msg: string, purpose: string): void {
        // console.info(msg);
        let dialogRef = this.dialog.open(DialogEditMachine, {
            //   width: '250px',
            data: { name: ip, message: msg, origin_purpose: purpose }
        });

        dialogRef.afterClosed().subscribe(result => {
            console.log('The dialog was closed');
        });
    }

    stopRunning(executor: string, executor_pid: string): void {
        let url_s: string = 'http://' + executor + ':9091/kill';
        let body: string = 'type=pid&id=' + executor_pid;
        this.http
            .post(url_s, body, {
                headers: new HttpHeaders().set('Content-Type', 'application/x-www-form-urlencoded'),
            }).subscribe(resp => {
                setTimeout(() => this.snackBar.open(('The kill request already send. IP:' + executor + '. PID:' + executor_pid), '', { duration: 5000 }));
            })
    }

}


export interface TestSummary {
    name: string;
    position: number;
    testindex: number;
    ci: string;
    status: string;
    result: string;
    resulturl: string;
    updatetime: string;
}


export interface RunningState {
    name: string;
    position: number;
    creator: string;
    progress: number;
    status: string;
    tools: string;
    createtime: string;
    color: string;
    create_url: string;
    status_info: string;
    executor: string;
    executor_pid: string;
}

export interface MachineState {
    name: string;
    position: number;
    info: string;
    update_time: string;
    test: string;
    status: string;
    message: string;
    purpose: string;
    message_update_time: string;
    labels: Array<string>;
    envinfos: Array<string>;
    uiinfo: string;
    uilink_status: string;
    screen_resolution: string;
}




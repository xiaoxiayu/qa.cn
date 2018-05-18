import { Component, OnInit, ViewChild, ChangeDetectorRef, Inject, ElementRef, OnDestroy } from '@angular/core';
import { HttpClientModule, HttpClient, HttpHeaders } from '@angular/common/http';
import { MatPaginator, MatTableDataSource, MatDialog, MatDialogRef, MAT_DIALOG_DATA } from '@angular/material';
import flvjs from 'flv.js';

@Component({
    selector: 'dialog-fuzz-live',
    templateUrl: 'dialog-fuzz-live.component.html',
    styleUrls: ['./dialog-fuzz-live.component.css']
})

export class DialogFuzzLive implements OnInit, OnDestroy {
    ip: string = "";
    name: string = "";
    loading: boolean = true;
    flvPlayer: any;

    @ViewChild("videoElement") myVideo: ElementRef;

    constructor(
        private http: HttpClient,
        public dialogRef: MatDialogRef<DialogFuzzLive>,
        @Inject(MAT_DIALOG_DATA) public data: any) {
        this.ip = data.ip;
        this.name = data.name;

    }

    ngOnDestroy(): void {
        // console.info('Destory');
        this.flvPlayer.pause();
        this.flvPlayer.unload();
        this.flvPlayer.detachMediaElement();
        this.flvPlayer.destroy();
        this.flvPlayer = null;
    }

    ngOnInit(): void {
        if (flvjs.isSupported()) {
            const videoElement = this.myVideo.nativeElement;
            console.info('http://' + this.ip + ':7001/live/test.flv');
            this.flvPlayer = flvjs.createPlayer({
                type: 'flv',
                url: 'http://' + this.ip + ':7001/live/test.flv'
            });
            this.flvPlayer.attachMediaElement(videoElement);
            this.loading = false;
            this.flvPlayer.load();
            this.flvPlayer.play();
            console.info('play start');
        }
    }

    onCancel(): void {
        this.dialogRef.close();
    }

    onOK(): void {

    }

    flv_start(): void {
        this.flvPlayer.play();
    }

    flv_pause(): void {
        this.flvPlayer.pause();
    }

}


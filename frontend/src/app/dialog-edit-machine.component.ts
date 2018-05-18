import { Component, OnInit, ViewChild, ChangeDetectorRef, Inject } from '@angular/core';
import { HttpClientModule, HttpClient, HttpHeaders } from '@angular/common/http';
import {MatPaginator, MatTableDataSource, MatDialog, MatDialogRef, MAT_DIALOG_DATA} from '@angular/material';

@Component({
    selector: 'dialog-edit-machine',
    templateUrl: 'dialog-edit-machine.component.html',
    styleUrls: ['./dialog-edit-machine.component.css']
  })
  
  export class DialogEditMachine {
    value: string = '';
    selectedValue: string;
    selected = '';

    constructor(
      private http: HttpClient,
      public dialogRef: MatDialogRef<DialogEditMachine>,
      @Inject(MAT_DIALOG_DATA) public data: any) { this.value = data.message; this.selected = data.origin_purpose;}
  
    onCancel(): void {
      this.dialogRef.close();
    }

    onOK(): void {
        if (this.value == this.data.message && this.selected == this.data.origin_purpose) {
            this.dialogRef.close();
            return;
        }
        let url_s: string = "test/state/gui-machine-msg";
        const body = "action=update&ip="+this.data.name+"&msg="+this.value+"&purpose="+this.selected;
        // console.info(body);
        this.http
            .post(url_s, body, {
                headers: new HttpHeaders().set('Content-Type', 'application/x-www-form-urlencoded'),
            }).subscribe(resp => {    
            console.info(resp);
            this.dialogRef.close();
        })
    }
  
  }

  
import { Component, OnInit, ViewChild, ChangeDetectorRef, Inject, ElementRef } from '@angular/core';
import { HttpClientModule, HttpClient, HttpHeaders } from '@angular/common/http';
import { MatPaginator, MatTableDataSource, MatDialog, MatDialogRef, MAT_DIALOG_DATA, MatSnackBar } from '@angular/material';
import { FileUploader, FileItem, ParsedResponseHeaders } from 'ng2-file-upload';

@Component({
  selector: 'dialog-uploadfiles',
  templateUrl: 'dialog-uploadfiles.component.html',
  styleUrls: ['./dialog-uploadfiles.component.css']
})

export class DialogUploadFiles {
  public upload_path: string = "";
  // file_service_url: string = 'http://127.0.0.1:8080/files/upload';
  file_service_url: string = 'files/upload';
  public uploader: FileUploader = new FileUploader({
    url: this.file_service_url,
    method: "POST",
    itemAlias: "file",
    autoUpload: false,
    isHTML5: true
  });

  infobackcolor: string = '';


  file_selected: boolean = false;
  selectedValue: string;
  selected = '';
  @ViewChild('upload_info') upload_info: ElementRef;

  constructor(
    private http: HttpClient,
    public snackBar: MatSnackBar,
    public dialogRef: MatDialogRef<DialogUploadFiles>,
    @Inject(MAT_DIALOG_DATA) public data: any
  ) {
    this.upload_path = this.urldecode(data.upload_to);
    // console.info(this.uppath);
    this.uploader.onAfterAddingFile = (file) => {
      file.withCredentials = false;
      if (file.file.size > 20971520) {
        file.remove();
        setTimeout(() => this.snackBar.open(('Can not upload "'+file.file.name+'" that size larger than 20M.'), '', { duration: 5000 }));
        return;
      }
      
    };

    this.uploader.onBuildItemForm = (fileItem: FileItem, form: any) => {
      form.append("FileName", fileItem.file.name);
      form.append("FileInfo", this.upload_info.nativeElement.value);
      form.append("FolderPath", this.upload_path);
    }

    this.uploader.onSuccessItem = (item: FileItem, response: string, status: number, headers: ParsedResponseHeaders) => {
      console.info('success');
    }

    // this.uploader.onCompleteItem = (item: FileItem, response: string, status: number, headers: ParsedResponseHeaders) => {
    //   console.info('complete');
    // }

    this.uploader.onSuccessItem = (item: FileItem, response: string, status: number, headers: ParsedResponseHeaders) => {

      let res = JSON.parse(response);
      if (res.ret == 2) {
        console.info(res);
        let fileid: string = res.info.split(' ')[0];
        fileid = fileid.split(':')[1];
        
        let url_s: string = 'files/_update_db_info';
        const body = "FileID=" + fileid + '&Info=' + this.upload_info.nativeElement.value;
        // console.info(body);
        this.http
          .post(url_s, body, {
            headers: new HttpHeaders().set('Content-Type', 'application/x-www-form-urlencoded'),
          }).subscribe(resp => {
            if (resp['affect'] == 1) {
              setTimeout(() => this.snackBar.open(('Info Updated: ' + item.file.name), '', { duration: 2000 }));
            } else {
              setTimeout(() => this.snackBar.open(('Info Updated Failed: ' + item.file.name), '', { duration: 10000 }));
            }
          })
      }
    }


    this.uploader.onErrorItem = (item: FileItem, response: string, status: number, headers: ParsedResponseHeaders) => {
      setTimeout(() => this.snackBar.open(('FAILED:' + item.file.name), '', { duration: 3000 }));

    }
  }

  doUpload(fileItem: FileItem): void {
    console.info(this.upload_info);
    if (this.upload_info.nativeElement.value == '') {
      this.infobackcolor = 'red';
      setTimeout(() => this.snackBar.open(('Info not set.'), '', { duration: 3000 }));
      setTimeout(() => { this.infobackcolor = ''; }, 1000);
      return;
    }
    fileItem.upload();
  }

  doUploadAll(): void {
    if (this.upload_info.nativeElement.value == '') {
      this.infobackcolor = 'red';
      setTimeout(() => this.snackBar.open(('Info not set.'), '', { duration: 3000 }));
      setTimeout(() => { this.infobackcolor = ''; }, 1000);
      return;
    }
    this.uploader.uploadAll();
  }

  urldecode(str: string): string {
    return decodeURIComponent((str + '').replace(/\+/g, '%20'));
  }

  selectedFileOnChanged(e: any): void {
    this.file_selected = true;
    console.info('select');
  }

  onClose(): void {
    this.dialogRef.close();
  }

  onReturn(): void {
    this.uploader.clearQueue();
    this.file_selected = false;
  }

  onOK(): void {
    // if (this.value == this.data.message && this.selected == this.data.origin_purpose) {
    //     this.dialogRef.close();
    //     return;
    // }
    // let url_s: string = "http://10.103.129.79/test/state/gui-machine-msg";
    // const body = "action=update&ip="+this.data.name+"&msg="+this.value+"&purpose="+this.selected;
    // // console.info(body);
    // this.http
    //     .post(url_s, body, {
    //         headers: new HttpHeaders().set('Content-Type', 'application/x-www-form-urlencoded'),
    //     }).subscribe(resp => {    
    //     console.info(resp);
    //     this.dialogRef.close();
    // })
  }

}
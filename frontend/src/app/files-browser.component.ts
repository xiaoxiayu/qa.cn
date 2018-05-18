import { Component, OnInit, ViewChild } from '@angular/core';
import { HttpClientModule, HttpClient } from '@angular/common/http';
import { FormControl } from '@angular/forms';
import { TreeModel, NodeEvent, NodeMenuItemAction, MenuItemSelectedEvent, Tree } from 'ng2-tree';
import { FileUploader } from 'ng2-file-upload';
import {
  MatDialog,
  MatDialogRef,
  MatSnackBar,
  MAT_DIALOG_DATA
} from '@angular/material';
import { DialogUploadFiles } from './dialog-uploadfiles.component';
import { WindowRef } from './fxqacommon.component';



@Component({
  selector: 'app-files-browser',
  templateUrl: './files-browser.component.html',
  // template: ``,
  styleUrls: ['./files-browser.component.css']
})

export class FilesBrowserComponent implements OnInit {
  nativeWindow: any;
  file_server_url: string = 'http://10.103.2.166:9090/TestFiles/';
  // file_server_url: string = 'http://127.0.0.1:9090/TestFiles';
  public tree: TreeModel = {
    settings: {
      isCollapsedOnInit: true,

      menuItems: [
        { action: NodeMenuItemAction.NewFolder, name: 'New Folder', cssClass: 'fa fa-arrow-right' },
        { action: NodeMenuItemAction.Custom, name: 'Upload', cssClass: 'fa fa-arrow-right' }
      ]
    },
    value: 'TestFiles/',
    id: '_root/',
    children: []
  };

  @ViewChild('treeComponent') treeComponent;

  constructor(
    private http: HttpClient,
    public snackBar: MatSnackBar,
    public dialog: MatDialog,
    private winRef: WindowRef) { 
      this.nativeWindow = winRef.nativeWindow;
    }

  public handleSelected(e: NodeEvent): void {
    if (e.node.foldingType.cssClass == 'node-leaf') {
      //setTimeout(() => this.snackBar.open(('Only folder can upload.'), '', { duration: 3000 }));
      let leaf_url: string = this.file_server_url + (<string>e.node.node.id).substring(6);
      this.nativeWindow.open(leaf_url);
      return;
    }
  }

  html_decode(str: string): string {
    var s = "";
    if (str.length == 0) return "";
    s = str.replace(/&gt;/g, "&");
    s = s.replace(/&lt;/g, "<");
    s = s.replace(/&gt;/g, ">");
    s = s.replace(/&nbsp;/g, " ");
    s = s.replace(/&#39;/g, "\'");
    s = s.replace(/&quot;/g, "\"");
    s = s.replace(/<br>/g, "\n");
    return s;
  }

  onMenuItemSelected(e: MenuItemSelectedEvent) {
    // console.log(e, `You selected ${e.selectedItem} menu item`);
    // console.log(e.node.parent);
    if (e.selectedItem == 'Upload') {
      if (e.node.foldingType.cssClass == 'node-leaf') {
        setTimeout(() => this.snackBar.open(('Only folder can use upload.'), '', { duration: 3000 }));
        return;
      }
      // console.info(e);
      let path_s: string = <string>e.node.node.id;
      if (path_s == '_root/') {
        let rootcontroler = this.treeComponent.getControllerByNodeId('_root/');
        this.openDialog('TestFiles/', rootcontroler);
        return;
      }
      path_s = path_s.substring(6);
      console.info('ID:', decodeURIComponent(path_s));
      console.info("Value:" + this.html_decode(<string>e.node.node.value));
      console.info('FoldingType:' + e.node.foldingType.cssClass);
      if (decodeURIComponent(path_s) != this.html_decode(<string>e.node.node.value)) {
        path_s = (<string>e.node.parent.id).substring(6) + this.html_decode(<string>e.node.node.value) + '/';
        // oopNodeController.collapse();
      }
      // console.info('child: ', e.node.parent.children);
      if (e.node.parent.children.length == 1) {
        if (e.node.parent.children[0].value == e.node.node.value) {
          setTimeout(() => this.snackBar.open(('Can not upload to "' + this.html_decode(<string>e.node.node.value) + '".' +
            ' Its parent "' + e.node.parent.value + '" is empty.'), '', { duration: 3000 }));
          return;
        }
      }
      let parentcontroler = this.treeComponent.getControllerByNodeId(e.node.parent.id);
      this.openDialog('TestFiles/' + path_s, parentcontroler);
    }

  }

  handleCollapsed(e: NodeEvent): void {
    // console.log(e);
  }

  handleCreated(e: NodeEvent): void {
    // console.info('fffffff');
  }

  handleExpanded(e: NodeEvent): void {
    let url_s: string = this.file_server_url + ((<string>e.node.node.id).substring(5)) + '?a=' + Date.now();
    console.info(url_s);
    this.http
      .get(url_s)
      .subscribe(
        (data) => { }, // Reach here if res.status >= 200 && <= 299
        (err) => {
          // console.info(err.error.text);
          let data = err.error.text;
          let data_list = data.split('\n');
          let newChildren: Array<TreeModel> = [];

          for (let i = 1; i < (data_list.length - 2); i++) {
            let l = data_list[i];
            let href_s: string = l.substring(l.indexOf('<a href="') + 9, l.indexOf('">'));
            l = l.replace('</a>', '');
            l = l.substring(l.indexOf('">') + 2);
            // console.info(l);
            // console.info(l, l.indexOf('/'), l.length);
            if (l.indexOf('/') != -1) {

              newChildren.push({
                settings: {
                  isCollapsedOnInit: true
                },
                value: l,
                id: e.node.node.id + href_s,
                children: [{
                  value: 'loading...'
                }]
              });
            } else {
              newChildren.push({
                value: l,
                id: e.node.node.id + '/' + href_s
              });
            }
          }
          const oopNodeController = this.treeComponent.getControllerByNodeId(e.node.node.id);
          oopNodeController.setChildren(newChildren);
        });
  }

  openDialog(path: string, parent_ctrl: any): void {
    // console.info(msg);
    let dialogRef = this.dialog.open(DialogUploadFiles, {
      //   width: '250px',
      height: '600px',
      width: '800px',
      data: { upload_to: path }
    });

    dialogRef.afterClosed().subscribe(result => {
      console.log('The dialog was closed');
      if (!parent_ctrl.isCollapsed()) {
        parent_ctrl.collapse();
        parent_ctrl.expand();
      }
    });
  }



  ngOnInit(): void {
    this.http
      .get(this.file_server_url)
      .subscribe(
        (data) => { }, // Reach here if res.status >= 200 && <= 299
        (err) => {
          // console.info(err.error.text);
          let data = err.error.text;
          let data_list = data.split('\n');
          let newChildren: Array<TreeModel> = [];

          for (let i = 1; i < (data_list.length - 2); i++) {
            let l = data_list[i];
            l = l.replace('</a>', '');

            l = l.substring(l.indexOf('">') + 2);
            // console.info(l, l.indexOf('/'), l.length);
            if (l.indexOf('/') != -1) {
              newChildren.push({ value: l, id: '_root/' + l, children: [{ value: 'loading...' }] });
            } else {
              newChildren.push({ value: l, id: '_root/' + l });
            }
          }
          const oopNodeController = this.treeComponent.getControllerByNodeId('_root/');
          oopNodeController.setChildren(newChildren);
        });

  }


}
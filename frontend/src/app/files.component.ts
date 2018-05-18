import { Component, OnInit, ViewChild } from '@angular/core';
import { HttpClientModule, HttpClient } from '@angular/common/http';
import { FormControl } from '@angular/forms';
import { DataTableDirective } from 'angular-datatables';
import { TreeModel, NodeEvent, NodeMenuItemAction, MenuItemSelectedEvent, Tree } from 'ng2-tree';
import { FileUploader } from 'ng2-file-upload';
import {
  MatDialog,
  MatDialogRef,
  MatSnackBar,
  MAT_DIALOG_DATA
} from '@angular/material';
import { DialogUploadFiles } from './dialog-uploadfiles.component';

class Person {
  id: number;
  firstName: string;
  lastName: string;
}

export class FileInfo {
  name: string;
  type: string;
  size: string;
  date: string;
  info: string;
  tools: string;
}

class DataTablesResponse {
  data: any[];
  draw: number;
  recordsFiltered: number;
  recordsTotal: number;
}


@Component({
  selector: 'app-files',
  templateUrl: './files.component.html',
  // template: ``,
  styleUrls: ['./files.component.css'],
  providers: [FileInfo]
})

export class FilesComponent implements OnInit {
  // file_server_url: string = 'http://10.103.2.166:9090/TestFiles/';
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


  @ViewChild(DataTableDirective)
  datatableElement: DataTableDirective;

  dtOptions: DataTables.Settings = {};
  fileinfos: FileInfo[];

  persons: Person[];
  lanOptions: DataTables.LanguageSettings = {};

  date_from = new FormControl('2015-02-12T14:18:43.488Z');
  date_to = new FormControl((new Date()).toISOString());

  constructor(
    private http: HttpClient,
    public searchdata: FileInfo,
    public dialog: MatDialog) { }
 

  simpleSearch(): void {
    this.datatableElement.dtInstance.then((dtInstance: DataTables.Api) => {
      if (this.searchdata.name == ''
        && this.searchdata.type == ''
        && this.searchdata.size == ''
        && this.searchdata.info == ''
        && this.date_from.value == '2015-02-12T14:18:43.488Z') {
        return;
      } else {
        console.info(this.searchdata);
      }
      // alert(this.date_from.value);
      // console.info(new Date(this.date_from.value).toISOString());
      //>= '2015-07-22' | <= '2017-01-01'
      this.searchdata.date = "Date >= '" +
        new Date(this.date_from.value).toLocaleDateString().split('T')[0].split('/').join('-')
        + "' AND Date <= '"
        + new Date(this.date_to.value).toLocaleDateString().split('T')[0].split('/').join('-')
        + "'";

      let simple_search = {
        FileName: this.searchdata.name,
        FileType: this.searchdata.type,
        Size: this.searchdata.size,
        Info: this.searchdata.info,
        Date: this.searchdata.date
      };

      let search_s: string = JSON.stringify(simple_search);
      console.info(search_s);
      dtInstance.search(search_s);
      dtInstance.draw();
    });
  }

  initCallBack(): void {
    const that = this;
    $("div.dataTables_filter input").unbind();
    // $('div.dataTables_filter label').append(' <button type="button" class="btn btn-default btn-sm" id="searchGo-complex" data-toggle="popover" data-placement="bottom" data-content="Advanced search.Under the form of a simple search is recommended" >Go</button>');          

  }

  ngOnInit(): void {
    this.searchdata.name = '';
    this.searchdata.type = '';
    this.searchdata.size = '';
    this.searchdata.info = '';
    const that = this;

    this.lanOptions = {
      zeroRecords: ""
    };

    this.dtOptions = {
      initComplete: this.initCallBack,
      ordering: false,
      searching: true,
      language: this.lanOptions,
      ajax: (dataTablesParameters: any, callback) => {
        // console.info('RequestData:' + JSON.stringify(dataTablesParameters));
        let url_s: string = 'files/getdata?draw=' +
          dataTablesParameters['draw'] +
          '&start=' + dataTablesParameters['start'] +
          '&length=' + dataTablesParameters['length'] +
          '&search[value]=' + encodeURIComponent(dataTablesParameters['search']['value']);
        that.http
          .get<DataTablesResponse>(url_s).subscribe(resp => {
            console.info('Respond:' + JSON.stringify(resp.data));
            that.fileinfos = resp.data;

            callback({
              recordsTotal: resp.recordsTotal,
              recordsFiltered: resp.recordsFiltered,
              data: []
            });
          });
      },
      processing: true,
      serverSide: true,
      columns: [{ data: 'name' }, { data: 'type' }, { data: 'size' }, { data: 'date' }, { data: 'info' }, { data: 'tools' }]
    };

  }

  tabChanged(event: any) {


  }
}
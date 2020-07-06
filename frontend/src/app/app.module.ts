import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { HttpModule } from '@angular/http';
// import { HttpClientModule, HttpClient } from '@angular/common/http';
import { HttpClientModule } from '@angular/common/http';

import { SafePipe, WindowRef } from './fxqacommon.component';
import { AppComponent } from './app.component';
import { CIComponent } from './ci.component';
import { FilesComponent } from './files.component';
import { StateComponent } from './state.component';
import { FuzzComponent } from './fuzz.component';
import { FilesBrowserComponent } from './files-browser.component'
import { DialogEditMachine } from './dialog-edit-machine.component';
import { DialogUploadFiles } from './dialog-uploadfiles.component';
import { DialogFuzzLive } from './dialog-fuzz-live.component';
import { ChartsComponent } from './charts.component';
import { xxJSComponent } from './xxjs.component';
import { AppRoutingModule } from './app-routing.module';

import { NgbModule } from '@ng-bootstrap/ng-bootstrap';
import { DataTablesModule } from 'angular-datatables';
import { TreeModule } from 'ng2-tree';
import { FileUploadModule } from 'ng2-file-upload';
import { NgxEchartsModule } from 'ngx-echarts';
import { AceEditorModule } from 'ng2-ace-editor';

import {
  MatAutocompleteModule,
  MatButtonModule,
  MatButtonToggleModule,
  MatCardModule,
  MatCheckboxModule,
  MatChipsModule,
  MatDatepickerModule,
  MatDialogModule,
  MatDividerModule,
  MatExpansionModule,
  MatGridListModule,
  MatIconModule,
  MatInputModule,
  MatListModule,
  MatMenuModule,
  MatNativeDateModule,
  MatPaginatorModule,
  MatProgressBarModule,
  MatProgressSpinnerModule,
  MatRadioModule,
  MatRippleModule,
  MatSelectModule,
  MatSidenavModule,
  MatSliderModule,
  MatSlideToggleModule,
  MatSnackBarModule,
  MatSortModule,
  MatStepperModule,
  MatTableModule,
  MatTabsModule,
  MatToolbarModule,
  MatTooltipModule,
} from '@angular/material';

import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { StateService } from './state.service';
import { CIService } from './ci.service';
import { PageNotFoundComponent } from './page-notfound.component';


@NgModule({
  declarations: [
    AppComponent,
    CIComponent,
    SafePipe,
    StateComponent,
    FilesComponent,
    PageNotFoundComponent,
    DialogEditMachine,
    DialogUploadFiles,
    DialogFuzzLive,
    FilesBrowserComponent,
    FuzzComponent,
    ChartsComponent,
    xxJSComponent
  ],
  imports: [
    BrowserModule,
    FormsModule,
    ReactiveFormsModule,
    AppRoutingModule,
    HttpClientModule,
    DataTablesModule,
    BrowserAnimationsModule,
    MatAutocompleteModule,
    MatButtonModule,
    MatButtonToggleModule,
    MatCardModule,
    MatCheckboxModule,
    MatChipsModule,
    MatDatepickerModule,
    MatDialogModule,
    MatDividerModule,
    MatExpansionModule,
    MatGridListModule,
    MatIconModule,
    MatInputModule,
    MatListModule,
    MatMenuModule,
    MatNativeDateModule,
    MatPaginatorModule,
    MatProgressBarModule,
    MatProgressSpinnerModule,
    MatRadioModule,
    MatRippleModule,
    MatSelectModule,
    MatSidenavModule,
    MatSliderModule,
    MatSlideToggleModule,
    MatSnackBarModule,
    MatSortModule,
    MatStepperModule,
    MatTableModule,
    MatTabsModule,
    MatToolbarModule,
    MatTooltipModule,
    NgbModule.forRoot(),
    TreeModule,
    FileUploadModule,
    NgxEchartsModule,
    AceEditorModule
  ],
  providers: [
    WindowRef,
    StateService,
    CIService
  ],
  entryComponents: [
    DialogEditMachine,
    DialogUploadFiles,
    DialogFuzzLive
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }

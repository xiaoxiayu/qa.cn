import { NgModule }             from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
 
import { CIComponent }   from './ci.component';
import { FilesComponent }   from './files.component';
import { StateComponent }   from './state.component';
import { PageNotFoundComponent } from './page-notfound.component';
import { FuzzComponent } from './fuzz.component';
import { FilesBrowserComponent } from './files-browser.component';
import { ChartsComponent } from './charts.component';
import { FoxitJSComponent } from './foxitjs.component';
 
const routes: Routes = [
 // { path: '', redirectTo: '/ci', pathMatch: 'full' },
  { path: 'ci',  component: CIComponent },
  { path: 'ci/:id',  component: CIComponent },
  { path: 'files-search',  component: FilesComponent },
  { path: 'files-browser',  component: FilesBrowserComponent },
  { path: 'test/state',  component: StateComponent },
  { path: 'test/state/:id',  component: StateComponent },
  { path: 'test/fuzz',  component: FuzzComponent },
  { path: 'chart',  component: ChartsComponent },
  { path: 'js',  component: FoxitJSComponent },
  { path: '**', component: PageNotFoundComponent },
];
 
@NgModule({
  imports: [ RouterModule.forRoot(routes) ],
  exports: [ RouterModule ]
})
export class AppRoutingModule {}
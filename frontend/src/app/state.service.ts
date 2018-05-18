import { Injectable } from '@angular/core';


@Injectable()
export class StateService {

    updateTab(id: string): Promise<string> {
        return Promise.resolve(id);
    }
}
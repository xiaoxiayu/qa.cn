import { Injectable } from '@angular/core';


@Injectable()
export class CIService {

    updateTab(id: string): Promise<string> {
        return Promise.resolve(id);
    }
}
/* tslint:disable:no-unused-variable */

import { TestBed, async, inject } from '@angular/core/testing';
import { FormInfoService } from './form-info.service';

describe('Service: FormInfo', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [FormInfoService]
    });
  });

  it('should ...', inject([FormInfoService], (service: FormInfoService) => {
    expect(service).toBeTruthy();
  }));
});

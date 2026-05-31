import { TestBed } from '@angular/core/testing';

import { Render } from './render';

describe('Render', () => {
  let service: Render;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(Render);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});

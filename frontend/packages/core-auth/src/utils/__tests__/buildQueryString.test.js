import { stringify } from 'qs';

import { buildQueryString } from '../buildQueryString';

jest.mock('qs', () => ({
  stringify: jest.fn(),
}));

describe('utils/buildQueryString', () => {
  it('should call qs.stringify', () => {
    buildQueryString({ foo: 'bar' });
    expect(stringify).toHaveBeenCalled();
  });
});

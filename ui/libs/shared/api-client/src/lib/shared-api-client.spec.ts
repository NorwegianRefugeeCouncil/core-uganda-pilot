import { sharedApiClient } from './shared-api-client';

describe('sharedApiClient', () => {
  it('should work', () => {
    expect(sharedApiClient()).toEqual('shared-api-client');
  });
});

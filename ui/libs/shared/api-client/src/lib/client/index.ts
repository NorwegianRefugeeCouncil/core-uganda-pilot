import { CoreClient, CoreInterface } from './core';
import { RESTClient } from '../shared-api-client';

export interface Interface {
  core(): CoreInterface
}

export class ClientSet implements Interface {
  public constructor(private restClient: RESTClient) {
  }

  core(): CoreInterface {
    return new CoreClient(this.restClient);
  }
}

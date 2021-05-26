import { CoreClient, CoreInterface } from './core';
import { RESTClient } from '../rest';
import { DiscoveryClient, DiscoveryInterface } from './discovery';

export interface Interface {
  core(): CoreInterface

  discovery(): DiscoveryInterface
}

export class ClientSet implements Interface {
  public constructor(private restClient: RESTClient) {
  }

  core(): CoreInterface {
    return new CoreClient(this.restClient);
  }

  discovery(): DiscoveryInterface {
    return new DiscoveryClient(this.restClient);
  }
}

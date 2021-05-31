import { RESTClient } from '../../rest';
import { DiscoveryV1Client, DiscoveryV1Interface } from './v1';

export interface DiscoveryInterface {
  v1(): DiscoveryV1Interface
}

export class DiscoveryClient implements DiscoveryInterface {
  public constructor(private client: RESTClient) {
  }

  v1(): DiscoveryV1Interface {
    return new DiscoveryV1Client(this.client);
  }
}

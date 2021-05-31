export * from './v1';

import { CoreV1Client, CoreV1Interface } from './v1';
import { RESTClient } from '../../rest';

export interface CoreInterface {
  v1(): CoreV1Interface
}

export class CoreClient {
  public constructor(private client: RESTClient) {
  }

  v1(): CoreV1Interface {
    return new CoreV1Client(this.client);
  }
}

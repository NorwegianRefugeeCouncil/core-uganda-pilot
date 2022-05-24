import { AdminClientDefinition } from '../types/client/AdminClient';

import { BaseRESTClient } from './BaseRESTClient';
import { IdentityClient } from './Identity';

export class Client extends BaseRESTClient implements AdminClientDefinition {
  static corev1 = 'apis/admin.nrc.no/v1';

  public Identity: IdentityClient;

  constructor(address: string) {
    super(`${address}/${Client.corev1}`);
    this.Identity = new IdentityClient(this);
  }
}

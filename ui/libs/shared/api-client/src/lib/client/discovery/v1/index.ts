import { RESTClient } from '../../../rest';
import { APIGroupsV1Client, APIGroupsV1Interface } from './apigroups';
import { APIResourcesV1Client, APIResourcesV1Interface } from './apiresources';

export interface DiscoveryV1Interface {
  apiGroups(): APIGroupsV1Interface

  apiResources(): APIResourcesV1Interface
}

export class DiscoveryV1Client implements DiscoveryV1Interface {
  public constructor(private restClient: RESTClient) {
  }

  apiGroups(): APIGroupsV1Interface {
    return new APIGroupsV1Client(this.restClient);
  }

  apiResources(): APIResourcesV1Interface {
    return new APIResourcesV1Client(this.restClient);
  }
}



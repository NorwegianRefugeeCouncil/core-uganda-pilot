import { RESTClient } from '../../../rest';
import { APIGroupList } from '../../../api';


export interface APIResourcesV1Interface {
  list(group: string, version: string): Promise<APIGroupList>
}

export class APIResourcesV1Client implements APIResourcesV1Interface {
  public constructor(private c: RESTClient) {
  }

  public list(group: string, version: string): Promise<APIGroupList> {
    return this.c.get()
      .group(group)
      .version(version)
      .do<APIGroupList>();
  }

}

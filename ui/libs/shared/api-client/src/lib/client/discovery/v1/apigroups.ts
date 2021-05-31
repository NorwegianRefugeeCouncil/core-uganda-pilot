import { RESTClient } from '../../../rest';
import { APIGroup } from '../../../api';


export interface APIGroupsV1Interface {
  get(name): Promise<APIGroup>
}

export class APIGroupsV1Client implements APIGroupsV1Interface {
  public constructor(private c: RESTClient) {
  }

  public get(name: string): Promise<APIGroup> {
    return this.c.get()
      .group(name)
      .do<APIGroup>();
  }

}

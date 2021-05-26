import { RESTClient } from '../../../rest';
import { APIGroupList } from '../../../api/discovery/v1';
import { Observable } from 'rxjs';

export interface APIResourcesV1Interface {
  list(group: string, version: string): Observable<APIGroupList>
}

export class APIResourcesV1Client implements APIResourcesV1Interface {
  public constructor(private c: RESTClient) {
  }

  public list(group: string, version: string): Observable<APIGroupList> {
    return this.c.get()
      .group(group)
      .version(version)
      .do<APIGroupList>();
  }

}

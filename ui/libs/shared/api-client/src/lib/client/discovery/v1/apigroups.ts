import { RESTClient } from '../../../rest';
import { APIGroup } from '../../../api/discovery/v1';
import { Observable } from 'rxjs';

export interface APIGroupsV1Interface {
  get(name): Observable<APIGroup>
}

export class APIGroupsV1Client implements APIGroupsV1Interface {
  public constructor(private c: RESTClient) {
  }

  public get(name: string): Observable<APIGroup> {
    return this.c.get()
      .group(name)
      .do<APIGroup>();
  }

}

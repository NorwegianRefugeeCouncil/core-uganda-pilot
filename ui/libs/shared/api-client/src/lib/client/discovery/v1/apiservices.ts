import { Observable } from 'rxjs';
import { APIServiceList } from '../../../api/discovery/v1';
import { RESTClient } from '../../../rest';

export interface APIServicesV1Interface {
  list(): Observable<APIServiceList>
}

export class APIServicesV1Client implements APIServicesV1Interface {
  public constructor(private c: RESTClient) {
  }

  public list(): Observable<APIServiceList> {
    return this.c.get()
      .group('discovery.nrc.no')
      .resource('apiservices')
      .do<APIServiceList>();
  }

}


import { APIServiceList } from '../../../api';
import { RESTClient } from '../../../rest';

export interface APIServicesV1Interface {
  list(): Promise<APIServiceList>
}

export class APIServicesV1Client implements APIServicesV1Interface {
  public constructor(private c: RESTClient) {
  }

  public list(): Promise<APIServiceList> {
    return this.c.get()
      .group('discovery.nrc.no')
      .resource('apiservices')
      .do<APIServiceList>();
  }

}

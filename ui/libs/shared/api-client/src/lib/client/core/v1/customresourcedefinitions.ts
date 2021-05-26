import { RESTClient } from '../../../rest';
import { CustomResourceDefinition, CustomResourceDefinitionList } from '../../../api/core/v1';
import { Observable } from 'rxjs';

export interface CustomResourceDefinitionsV1Interface {
  get(name): Observable<CustomResourceDefinition>

  create(customResourceDefinition: CustomResourceDefinition): Observable<CustomResourceDefinition>

  update(customResourceDefinition: CustomResourceDefinition): Observable<CustomResourceDefinition>

  delete(name: string): Observable<void>

  list(): Observable<CustomResourceDefinitionList>
}


export class CustomResourceDefinitionsV1Client implements CustomResourceDefinitionsV1Interface {
  public constructor(private c: RESTClient) {
  }

  public get(name: string): Observable<CustomResourceDefinition> {
    return this.c.get()
      .version('v1')
      .resource('customresourcedefinitions')
      .group('core.nrc.no')
      .name(name)
      .do<CustomResourceDefinition>();
  }

  public create(formDefinition: CustomResourceDefinition): Observable<CustomResourceDefinition> {
    return this.c.post()
      .group('core.nrc.no')
      .version('v1')
      .resource('customresourcedefinitions')
      .body(formDefinition)
      .do<CustomResourceDefinition>();
  }

  public update(formDefinition: CustomResourceDefinition): Observable<CustomResourceDefinition> {
    return this.c.put()
      .group('core.nrc.no')
      .version('v1')
      .resource('customresourcedefinitions')
      .name(formDefinition.metadata.name)
      .body(formDefinition)
      .do<CustomResourceDefinition>();
  }

  public delete(name: string): Observable<void> {
    return this.c.delete()
      .group('core.nrc.no')
      .version('v1')
      .resource('customresourcedefinitions')
      .name(name)
      .do<void>();
  }

  public list(): Observable<CustomResourceDefinitionList> {
    return this.c.get()
      .version('v1')
      .resource('customresourcedefinitions')
      .group('core.nrc.no')
      .do<CustomResourceDefinitionList>();
  }

}

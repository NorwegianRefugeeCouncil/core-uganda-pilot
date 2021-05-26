import { RESTClient } from '../../../rest';
import { FormDefinition, FormDefinitionList } from '../../../api/core/v1';
import { Observable } from 'rxjs';

export interface FormDefinitionsV1Interface {
  get(name): Observable<FormDefinition>

  create(formDefinition: FormDefinition): Observable<FormDefinition>

  update(formDefinition: FormDefinition): Observable<FormDefinition>

  delete(name: string): Observable<void>

  list(): Observable<FormDefinitionList>
}

export class FormDefinitionsV1Client implements FormDefinitionsV1Interface {
  public constructor(private c: RESTClient) {
  }

  public get(name: string): Observable<FormDefinition> {
    return this.c.get()
      .version('v1')
      .resource('formdefinitions')
      .group('core.nrc.no')
      .name(name)
      .do<FormDefinition>();
  }

  public create(formDefinition: FormDefinition): Observable<FormDefinition> {
    return this.c.post()
      .group('core.nrc.no')
      .version('v1')
      .resource('formdefinitions')
      .body(formDefinition)
      .do<FormDefinition>();
  }

  public update(formDefinition: FormDefinition): Observable<FormDefinition> {
    return this.c.put()
      .group('core.nrc.no')
      .version('v1')
      .resource('formdefinitions')
      .name(formDefinition.metadata.name)
      .body(formDefinition)
      .do<FormDefinition>();
  }

  public delete(name: string): Observable<void> {
    return this.c.delete()
      .group('core.nrc.no')
      .version('v1')
      .resource('formdefinitions')
      .name(name)
      .do<void>();
  }

  public list(): Observable<FormDefinitionList> {
    return this.c.get()
      .version('v1')
      .resource('formdefinitions')
      .group('core.nrc.no')
      .do<FormDefinitionList>();
  }

}

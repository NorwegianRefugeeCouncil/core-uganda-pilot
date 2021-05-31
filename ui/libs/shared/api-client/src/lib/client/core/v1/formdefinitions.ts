import { RESTClient } from '../../../rest';
import { FormDefinition, FormDefinitionList } from '../../../api';


export interface FormDefinitionsV1Interface {
  get(name): Promise<FormDefinition>

  create(formDefinition: FormDefinition): Promise<FormDefinition>

  update(formDefinition: FormDefinition): Promise<FormDefinition>

  delete(name: string): Promise<void>

  list(): Promise<FormDefinitionList>
}

export class FormDefinitionsV1Client implements FormDefinitionsV1Interface {
  public constructor(private c: RESTClient) {
  }

  public get(name: string): Promise<FormDefinition> {
    return this.c.get()
      .version('v1')
      .resource('formdefinitions')
      .group('core.nrc.no')
      .name(name)
      .do<FormDefinition>();
  }

  public create(formDefinition: FormDefinition): Promise<FormDefinition> {
    return this.c.post()
      .group('core.nrc.no')
      .version('v1')
      .resource('formdefinitions')
      .body(formDefinition)
      .do<FormDefinition>();
  }

  public update(formDefinition: FormDefinition): Promise<FormDefinition> {
    return this.c.put()
      .group('core.nrc.no')
      .version('v1')
      .resource('formdefinitions')
      .name(formDefinition.metadata.name)
      .body(formDefinition)
      .do<FormDefinition>();
  }

  public delete(name: string): Promise<void> {
    return this.c.delete()
      .group('core.nrc.no')
      .version('v1')
      .resource('formdefinitions')
      .name(name)
      .do<void>();
  }

  public list(): Promise<FormDefinitionList> {
    return this.c.get()
      .version('v1')
      .resource('formdefinitions')
      .group('core.nrc.no')
      .do<FormDefinitionList>();
  }

}

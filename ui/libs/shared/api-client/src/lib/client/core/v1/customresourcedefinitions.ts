import { RESTClient } from '../../../rest';
import { CustomResourceDefinition, CustomResourceDefinitionList } from '../../../api';


export interface CustomResourceDefinitionsV1Interface {
  get(name): Promise<CustomResourceDefinition>

  create(customResourceDefinition: CustomResourceDefinition): Promise<CustomResourceDefinition>

  update(customResourceDefinition: CustomResourceDefinition): Promise<CustomResourceDefinition>

  delete(name: string): Promise<void>

  list(): Promise<CustomResourceDefinitionList>
}


export class CustomResourceDefinitionsV1Client implements CustomResourceDefinitionsV1Interface {
  public constructor(private c: RESTClient) {
  }

  public get(name: string): Promise<CustomResourceDefinition> {
    return this.c.get()
      .version('v1')
      .resource('customresourcedefinitions')
      .group('core.nrc.no')
      .name(name)
      .do<CustomResourceDefinition>();
  }

  public create(formDefinition: CustomResourceDefinition): Promise<CustomResourceDefinition> {
    return this.c.post()
      .group('core.nrc.no')
      .version('v1')
      .resource('customresourcedefinitions')
      .body(formDefinition)
      .do<CustomResourceDefinition>();
  }

  public update(formDefinition: CustomResourceDefinition): Promise<CustomResourceDefinition> {
    return this.c.put()
      .group('core.nrc.no')
      .version('v1')
      .resource('customresourcedefinitions')
      .name(formDefinition.metadata.name)
      .body(formDefinition)
      .do<CustomResourceDefinition>();
  }

  public delete(name: string): Promise<void> {
    return this.c.delete()
      .group('core.nrc.no')
      .version('v1')
      .resource('customresourcedefinitions')
      .name(name)
      .do<void>();
  }

  public list(): Promise<CustomResourceDefinitionList> {
    return this.c.get()
      .version('v1')
      .resource('customresourcedefinitions')
      .group('core.nrc.no')
      .do<CustomResourceDefinitionList>();
  }

}

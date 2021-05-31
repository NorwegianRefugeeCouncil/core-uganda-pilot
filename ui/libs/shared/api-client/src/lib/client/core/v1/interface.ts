import { FormDefinitionsV1Client, FormDefinitionsV1Interface } from './formdefinitions';
import { CustomResourceDefinitionsV1Client, CustomResourceDefinitionsV1Interface } from './customresourcedefinitions';
import { RESTClient } from '../../../rest';

export interface CoreV1Interface {
  formDefinitions(): FormDefinitionsV1Interface

  customResourceDefinitions(): CustomResourceDefinitionsV1Interface
}

export class CoreV1Client implements CoreV1Interface {
  public constructor(private restClient: RESTClient) {
  }

  formDefinitions(): FormDefinitionsV1Interface {
    return new FormDefinitionsV1Client(this.restClient);
  }

  customResourceDefinitions(): CustomResourceDefinitionsV1Interface {
    return new CustomResourceDefinitionsV1Client(this.restClient);
  }
}



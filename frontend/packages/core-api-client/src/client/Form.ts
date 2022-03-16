import {
  FormCreateRequest,
  FormCreateResponse,
  FormGetRequest,
  FormGetResponse,
  FormListResponse,
  FormDefinition,
  FieldDefinition,
} from '../types';

import { BaseRESTClient } from './BaseRESTClient';

export class FormClient {
  restClient: BaseRESTClient;

  constructor(restClient: BaseRESTClient) {
    this.restClient = restClient;
  }

  create = (request: FormCreateRequest): Promise<FormCreateResponse> => {
    return this.restClient.do(request, '/forms', 'post', request.object, 200);
  };

  list = (request: {} | undefined): Promise<FormListResponse> => {
    return this.restClient.do(request, '/forms', 'get', undefined, 200);
  };

  get = (request: FormGetRequest): Promise<FormGetResponse> => {
    return this.restClient.do(
      request,
      `/forms/${request.id}`,
      'get',
      undefined,
      200,
    );
  };

  getAncestors = async (formId: string): Promise<FormDefinition[]> => {
    const formResponse = await this.get({ id: formId });

    if (!formResponse.response) {
      throw new Error(formResponse.error);
    }

    const referenceKey = formResponse.response.fields.find(
      (field: FieldDefinition) => field.key && field.fieldType.reference,
    );

    if (referenceKey && referenceKey.fieldType.reference) {
      const result = await this.getAncestors(
        referenceKey.fieldType.reference.formId,
      );
      return [...result, formResponse.response];
    }
    return [formResponse.response];
  };
}

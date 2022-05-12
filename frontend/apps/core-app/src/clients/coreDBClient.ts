/*
  USED FOR TESTING PURPOSES ONLY
*/

import { BaseRESTClient } from 'core-api-client';

type EntityAttribute = {
  id: string;
  name: string;
  list: boolean;
  type:
    | 'string'
    | 'number'
    | 'boolean'
    | 'date'
    | 'time'
    | 'datetime'
    | 'month'
    | 'week'
    | 'co-ordinate'
    | 'file';
  constraints: {
    required?: boolean;
    unique?: boolean;
    min?: number;
    max?: number;
    pattern?: string;
    enum?: string[];
    custom?: string[];
  };
};

type EntityRelationship = {
  id: string;
  cardinality: 'one-to-one' | 'one-to-many' | 'many-to-one' | 'many-to-many';
  sourceEntityId: string;
  targetEntityId: string;
};

type Entity = {
  id: string;
  name: string;
  description: string;
  attributes: EntityAttribute[];
  constraints: {
    custom?: string[];
  };
  relationships: EntityRelationship[];
};

export const dummyEntity: Entity = {
  id: 'e1',
  name: 'My Entity',
  description: '',
  attributes: [
    {
      id: 'a1',
      name: 'My Attribute',
      list: false,
      type: 'string',
      constraints: {
        required: true,
        min: 0,
        max: 10,
      },
    },
  ],
  constraints: {
    custom: [],
  },
  relationships: [],
};

export class Client extends BaseRESTClient {
  static corev1 = 'apis/core.nrc.no/v1';

  constructor(address: string) {
    super(`${address}/${Client.corev1}`);
  }

  public createEntity = async (entity: Entity): Promise<Entity> => {
    const response = await this.post<Entity, Entity>('/entity', entity);
    if (response.error || !response.response) {
      throw new Error(response.error);
    }
    return response.response;
  };
}

export const coreDBClient = new Client('https://localhost:9010');

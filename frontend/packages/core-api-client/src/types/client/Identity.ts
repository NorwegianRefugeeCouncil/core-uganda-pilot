import { DataOperation, Response } from './utils';

export type Identity = {
  id: string;
  subject: string;
  displayName: string;
  fullName: string;
  email: string;
  emailVerified: boolean;
};

export type IdentityGetRequest = { id: string };

export type IdentityGetResponse = Response<undefined, Identity>;

export interface IdentityClientDefinition {
  get: DataOperation<IdentityGetRequest, IdentityGetResponse>;
}

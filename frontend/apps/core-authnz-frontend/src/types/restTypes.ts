import {
  IdentityProvider,
  IdentityProviderList,
  OAuth2Client,
  OAuth2ClientList,
  Organization,
  OrganizationList,
  Session,
} from './types';
import { DataOperation, Response } from './utils';

export type SessionGetRequest = void;
export type SessionGetResponse = Response<SessionGetRequest, Session>;
export type SessionGetter = { getSession: DataOperation<SessionGetRequest, SessionGetResponse> };
export type OrganizationListRequest = void;
export type OrganizationListResponse = Response<OrganizationListRequest, OrganizationList>;
export type OrganizationLister = {
  listOrganizations: DataOperation<OrganizationListRequest, OrganizationListResponse>;
};
export type OrganizationCreateRequest = { object: Partial<Organization> };
export type OrganizationCreateResponse = Response<OrganizationCreateRequest, Organization>;
export type OrganizationCreator = {
  createOrganization: DataOperation<OrganizationCreateRequest, OrganizationCreateResponse>;
};
export type OrganizationGetRequest = { id: string };
export type OrganizationGetResponse = Response<OrganizationGetRequest, Organization>;
export type OrganizationGetter = {
  getOrganization: DataOperation<OrganizationGetRequest, OrganizationGetResponse>;
};
export type IdentityProviderGetRequest = { id: string };
export type IdentityProviderGetResponse = Response<IdentityProviderGetRequest, IdentityProvider>;
export type IdentityProviderGetter = {
  getIdentityProvider: DataOperation<IdentityProviderGetRequest, IdentityProviderGetResponse>;
};
export type IdentityProviderListRequest = { organizationId: string };
export type IdentityProviderListResponse = Response<IdentityProviderListRequest, IdentityProviderList>;
export type IdentityProviderLister = {
  listIdentityProviders: DataOperation<IdentityProviderListRequest, IdentityProviderListResponse>;
};
export type IdentityProviderCreateRequest = { object: Partial<IdentityProvider> };
export type IdentityProviderCreateResponse = Response<IdentityProviderCreateRequest, IdentityProvider>;
export type IdentityProviderCreator = {
  createIdentityProvider: DataOperation<IdentityProviderCreateRequest, IdentityProviderCreateResponse>;
};
export type IdentityProviderUpdateRequest = { object: Partial<IdentityProvider> };
export type IdentityProviderUpdateResponse = Response<IdentityProviderUpdateRequest, IdentityProvider>;
export type IdentityProviderUpdater = {
  updateIdentityProvider: DataOperation<IdentityProviderUpdateRequest, IdentityProviderUpdateResponse>;
};
export type OAuth2ClientListRequest = {};
export type OAuth2ClientListResponse = Response<OAuth2ClientListRequest, OAuth2ClientList>;
export type OAuth2ClientLister = {
  listOAuth2Clients: DataOperation<OAuth2ClientListRequest, OAuth2ClientListResponse>;
};
export type OAuth2ClientGetRequest = { id: string };
export type OAuth2ClientGetResponse = Response<OAuth2ClientGetRequest, OAuth2Client>;
export type OAuth2ClientGetter = {
  getOAuth2Client: DataOperation<OAuth2ClientGetRequest, OAuth2ClientGetResponse>;
};
export type OAuth2ClientUpdateRequest = { object: Partial<OAuth2Client> };
export type OAuth2ClientUpdateResponse = Response<OAuth2ClientUpdateRequest, OAuth2Client>;
export type OAuth2ClientUpdater = {
  updateOAuth2Client: DataOperation<OAuth2ClientUpdateRequest, OAuth2ClientUpdateResponse>;
};
export type OAuth2ClientCreateRequest = { object: Partial<OAuth2Client> };
export type OAuth2ClientCreateResponse = Response<OAuth2ClientCreateRequest, OAuth2Client>;
export type OAuth2ClientCreator = {
  createOAuth2Client: DataOperation<OAuth2ClientCreateRequest, OAuth2ClientCreateResponse>;
};
export type OAuth2ClientDeleteRequest = { id: string };
export type OAuth2ClientDeleteResponse = Response<OAuth2ClientDeleteRequest, void>;
export type OAuth2ClientDeleter = {
  deleteOAuth2Client: DataOperation<OAuth2ClientDeleteRequest, OAuth2ClientDeleteResponse>;
};

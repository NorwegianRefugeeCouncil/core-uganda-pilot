import { BaseRESTClient } from 'core-api-client';

import { Client as ApiClient } from '../types/types';
import {
  IdentityProviderCreateRequest,
  IdentityProviderCreateResponse,
  IdentityProviderGetRequest,
  IdentityProviderGetResponse,
  IdentityProviderListRequest,
  IdentityProviderListResponse,
  OAuth2ClientCreateRequest,
  OAuth2ClientCreateResponse,
  OAuth2ClientDeleteRequest,
  OAuth2ClientDeleteResponse,
  OAuth2ClientGetRequest,
  OAuth2ClientGetResponse,
  OAuth2ClientListRequest,
  OAuth2ClientListResponse,
  OAuth2ClientUpdateRequest,
  OAuth2ClientUpdateResponse,
  OrganizationCreateRequest,
  OrganizationCreateResponse,
  OrganizationGetResponse,
  OrganizationListResponse,
  SessionGetResponse,
} from '../types/restTypes';

export default class Client extends BaseRESTClient implements ApiClient {
  static adminv1 = 'apis/admin.nrc.no/v1';

  public constructor(address: string) {
    super(`${address}/${Client.adminv1}`);
  }

  createIdentityProvider(request: IdentityProviderCreateRequest): Promise<IdentityProviderCreateResponse> {
    return this.do(request, '/identityproviders', 'post', request.object, 200);
  }

  createOrganization(request: OrganizationCreateRequest): Promise<OrganizationCreateResponse> {
    return this.do(request, '/organizations', 'post', request.object, 200);
  }

  createOAuth2Client(request: OAuth2ClientCreateRequest): Promise<OAuth2ClientCreateResponse> {
    return this.do(request, '/clients', 'post', request.object, 200);
  }

  deleteOAuth2Client(request: OAuth2ClientDeleteRequest): Promise<OAuth2ClientDeleteResponse> {
    return this.do(request, `/clients/${request.id}`, 'delete', undefined, 204);
  }

  getIdentityProvider(request: IdentityProviderGetRequest): Promise<IdentityProviderGetResponse> {
    return this.do(request, `/identityproviders/${request.id}`, 'get', undefined, 200);
  }

  getOrganization(request: { id: string }): Promise<OrganizationGetResponse> {
    return this.do(request, `/organizations/${request.id}`, 'get', undefined, 200);
  }

  getOAuth2Client(request: OAuth2ClientGetRequest): Promise<OAuth2ClientGetResponse> {
    return this.do(request, `/clients/${request.id}`, 'get', undefined, 200);
  }

  listIdentityProviders(request: IdentityProviderListRequest): Promise<IdentityProviderListResponse> {
    return this.do(request, `/identityproviders?organizationId=${request.organizationId}`, 'get', undefined, 200);
  }

  listOAuth2Clients(request: OAuth2ClientListRequest): Promise<OAuth2ClientListResponse> {
    return this.do(request, '/clients', 'get', undefined, 200);
  }

  listOrganizations(request: void): Promise<OrganizationListResponse> {
    return this.do(request, '/organizations', 'get', undefined, 200);
  }

  updateIdentityProvider(request: IdentityProviderCreateRequest): Promise<IdentityProviderCreateResponse> {
    return this.do(request, `/identityproviders/${request.object.id}`, 'put', request.object, 200);
  }

  updateOAuth2Client(request: OAuth2ClientUpdateRequest): Promise<OAuth2ClientUpdateResponse> {
    return this.do(request, `/clients/${request.object.id}`, 'put', request.object, 200);
  }

  getSession(request: void): Promise<SessionGetResponse> {
    return this.do(request, '/oidc/session', 'get', undefined, 200, {
      headers: {},
    });
  }
}

import axios, {AxiosError, AxiosInstance, Method} from "axios";
import {Client as ApiClient} from '../types/types';
import {clientResponse} from "../utils/responses";
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
    SessionGetResponse
} from "../types/restTypes";
import {RequestOptions, Response} from "../types/utils";

export default class Client implements ApiClient {
    private adminv1 = 'apis/admin.nrc.no/v1';

    public constructor(
        public readonly address = 'https://localhost:9001/',
        public readonly axiosInstance: AxiosInstance = axios.create()) {
    }

    do<TRequest, TBody>(request: TRequest, url: string, method: Method, data: any, expectStatusCode: number, options?: RequestOptions): Promise<Response<TRequest, TBody>> {

        let headers: { [key: string]: string } = {
            "Accept": "application/json",
        }
        if (options?.headers) {
            headers = options?.headers
        }

        return axios.request<TBody>({
            method,
            url,
            data,
            responseType: "json",
            headers,
            withCredentials: true,
        }).then(value => {
            return clientResponse<TRequest, TBody>(value, request, expectStatusCode);
        }).catch((err: AxiosError) => {
            return {
                request: request,
                response: undefined,
                status: "500 Internal Server Error",
                statusCode: 500,
                error: err.message,
                success: false,
            }
        })
    }

    createIdentityProvider(request: IdentityProviderCreateRequest): Promise<IdentityProviderCreateResponse> {
        return this.do(request, `${this.address}/${this.adminv1}/identityproviders`, "post", request.object, 200)
    }

    createOrganization(request: OrganizationCreateRequest): Promise<OrganizationCreateResponse> {
        return this.do(request, `${this.address}/${this.adminv1}/organizations`, "post", request.object, 200)
    }

    createOAuth2Client(request: OAuth2ClientCreateRequest): Promise<OAuth2ClientCreateResponse> {
        return this.do(request, `${this.address}/${this.adminv1}/clients`, "post", request.object, 200)
    }

    deleteOAuth2Client(request: OAuth2ClientDeleteRequest): Promise<OAuth2ClientDeleteResponse> {
        return this.do(request, `${this.address}/${this.adminv1}/clients/${request.id}`, "delete", undefined, 204)
    }

    getIdentityProvider(request: IdentityProviderGetRequest): Promise<IdentityProviderGetResponse> {
        return this.do(request, `${this.address}/${this.adminv1}/identityproviders/${request.id}`, "get", undefined, 200)
    }

    getOrganization(request: { id: string }): Promise<OrganizationGetResponse> {
        return this.do(request, `${this.address}/${this.adminv1}/organizations/${request.id}`, "get", undefined, 200)
    }

    getOAuth2Client(request: OAuth2ClientGetRequest): Promise<OAuth2ClientGetResponse> {
        return this.do(request, `${this.address}/${this.adminv1}/clients/${request.id}`, "get", undefined, 200)
    }

    listIdentityProviders(request: IdentityProviderListRequest): Promise<IdentityProviderListResponse> {
        return this.do(request, `${this.address}/${this.adminv1}/identityproviders?organizationId=${request.organizationId}`, "get", undefined, 200)
    }

    listOAuth2Clients(request: OAuth2ClientListRequest): Promise<OAuth2ClientListResponse> {
        return this.do(request, `${this.address}/${this.adminv1}/clients`, "get", undefined, 200)
    }

    listOrganizations(request: void): Promise<OrganizationListResponse> {
        return this.do(request, `${this.address}/${this.adminv1}/organizations`, "get", undefined, 200)
    }

    updateIdentityProvider(request: IdentityProviderCreateRequest): Promise<IdentityProviderCreateResponse> {
        return this.do(request, `${this.address}/${this.adminv1}/identityproviders/${request.object.id}`, "put", request.object, 200)
    }

    updateOAuth2Client(request: OAuth2ClientUpdateRequest): Promise<OAuth2ClientUpdateResponse> {
        return this.do(request, `${this.address}/${this.adminv1}/clients/${request.object.id}`, "put", request.object, 200)
    }

    getSession(request: void): Promise<SessionGetResponse> {
        return this.do(request, `${this.address}/${this.adminv1}/oidc/session`, "get", undefined, 200, {headers: {}})
    }
}


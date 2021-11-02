import axios, {AxiosResponse, Method} from "axios";

export type DataOperation<TRequest, TResponse> = (request: TRequest) => Promise<TResponse>

export type Response<TRequest, TResponse> = {
    request: TRequest
    response: TResponse | undefined
    status: string
    statusCode: number
    success: boolean
    error: any
}

export class Organization {
    public id: string = ""
    public key: string = ""
    public name: string = ""
}

export type OrganizationList = { items: Organization[] }

export class IdentityProvider {
    public id: string = ""
    public name: string = ""
    public organizationId: string = ""
    public domain: string = ""
    public clientId: string = ""
    public clientSecret: string = ""
    public emailDomain: string = ""
}

export type IdentityProviderList = { items: IdentityProvider[] }

export type TokenEndpointAuthMethod = "client_secret_post" | "client_secret_basic" | "private_key_jwt" | "none"

export type ResponseType = "code" | "token" | "id_token"

export type GrantType = "authorization_code" | "refresh_token" | "client_credentials" | "implicit"

export class OAuth2Client {
    public id: string = ""
    public clientName: string = ""
    public clientSecret: string = ""
    public uri: string = ""
    public grantTypes: GrantType[] = ["authorization_code"]
    public responseTypes: ResponseType[] = ["code"]
    public scope: string = ""
    public redirectUris: string[] = []
    public allowedCorsOrigins: string[] = []
    public tokenEndpointAuthMethod: TokenEndpointAuthMethod = "client_secret_basic"
}

export type OAuth2ClientList = {
    items: OAuth2Client[]
}

export type OrganizationListRequest = void
export type OrganizationListResponse = Response<OrganizationListRequest, OrganizationList>
export type OrganizationLister = { listOrganizations: DataOperation<OrganizationListRequest, OrganizationListResponse> }
export type OrganizationCreateRequest = { object: Partial<Organization> }
export type OrganizationCreateResponse = Response<OrganizationCreateRequest, Organization>
export type OrganizationCreator = { createOrganization: DataOperation<OrganizationCreateRequest, OrganizationCreateResponse> }
export type OrganizationGetRequest = { id: string }
export type OrganizationGetResponse = Response<OrganizationGetRequest, Organization>
export type OrganizationGetter = { getOrganization: DataOperation<OrganizationGetRequest, OrganizationGetResponse> }
export type IdentityProviderGetRequest = { id: string }
export type IdentityProviderGetResponse = Response<IdentityProviderGetRequest, IdentityProvider>
export type IdentityProviderGetter = { getIdentityProvider: DataOperation<IdentityProviderGetRequest, IdentityProviderGetResponse> }
export type IdentityProviderListRequest = { organizationId: string }
export type IdentityProviderListResponse = Response<IdentityProviderListRequest, IdentityProviderList>
export type IdentityProviderLister = { listIdentityProviders: DataOperation<IdentityProviderListRequest, IdentityProviderListResponse> }
export type IdentityProviderCreateRequest = { object: Partial<IdentityProvider> }
export type IdentityProviderCreateResponse = Response<IdentityProviderCreateRequest, IdentityProvider>
export type IdentityProviderCreator = { createIdentityProvider: DataOperation<IdentityProviderCreateRequest, IdentityProviderCreateResponse> }
export type IdentityProviderUpdateRequest = { object: Partial<IdentityProvider> }
export type IdentityProviderUpdateResponse = Response<IdentityProviderUpdateRequest, IdentityProvider>
export type IdentityProviderUpdater = { updateIdentityProvider: DataOperation<IdentityProviderUpdateRequest, IdentityProviderUpdateResponse> }
export type OAuth2ClientListRequest = {}
export type OAuth2ClientListResponse = Response<OAuth2ClientListRequest, OAuth2ClientList>
export type OAuth2ClientLister = { listOAuth2Clients: DataOperation<OAuth2ClientListRequest, OAuth2ClientListResponse> }
export type OAuth2ClientGetRequest = { id: string }
export type OAuth2ClientGetResponse = Response<OAuth2ClientGetRequest, OAuth2Client>
export type OAuth2ClientGetter = { getOAuth2Client: DataOperation<OAuth2ClientGetRequest, OAuth2ClientGetResponse> }
export type OAuth2ClientUpdateRequest = { object: Partial<OAuth2Client> }
export type OAuth2ClientUpdateResponse = Response<OAuth2ClientUpdateRequest, OAuth2Client>
export type OAuth2ClientUpdater = { updateOAuth2Client: DataOperation<OAuth2ClientUpdateRequest, OAuth2ClientUpdateResponse> }
export type OAuth2ClientCreateRequest = { object: Partial<OAuth2Client> }
export type OAuth2ClientCreateResponse = Response<OAuth2ClientCreateRequest, OAuth2Client>
export type OAuth2ClientCreator = { createOAuth2Client: DataOperation<OAuth2ClientCreateRequest, OAuth2ClientCreateResponse> }
export type OAuth2ClientDeleteRequest = { id: string }
export type OAuth2ClientDeleteResponse = Response<OAuth2ClientDeleteRequest, void>
export type OAuth2ClientDeleter = { deleteOAuth2Client: DataOperation<OAuth2ClientDeleteRequest, OAuth2ClientDeleteResponse> }


export interface Client
    extends OrganizationLister,
        OrganizationCreator,
        OrganizationGetter,
        IdentityProviderGetter,
        IdentityProviderLister,
        IdentityProviderCreator,
        IdentityProviderUpdater,
        OAuth2ClientGetter,
        OAuth2ClientLister,
        OAuth2ClientUpdater,
        OAuth2ClientCreator,
        OAuth2ClientDeleter {
}

function errorResponse<TRequest, TBody>(request: TRequest, r: AxiosResponse<TBody>): Response<TRequest, TBody> {
    return {
        request: request,
        response: undefined,
        status: r.request,
        statusCode: r.status,
        error: r.data as any,
        success: false,
    };
}

function successResponse<TRequest, TBody>(request: TRequest, r: AxiosResponse<TBody>): Response<TRequest, TBody> {
    return {
        request: request,
        response: r.data as TBody,
        status: r.statusText,
        statusCode: r.status,
        error: undefined,
        success: true,
    };
}

function clientResponse<TRequest, TBody>(r: AxiosResponse<TBody>, request: TRequest, expectedStatusCode: number): Response<TRequest, TBody> {
    return r.status !== expectedStatusCode
        ? errorResponse<TRequest, TBody>(request, r)
        : successResponse<TRequest, TBody>(request, r)
}

export type clientProps = {
    idToken?: string
    address?: string
}

export class client implements Client {
    public address = "http://localhost:9001/admin"
    public idToken = ""

    public constructor(private clientProps?: clientProps) {
        if (clientProps?.idToken) {
            this.idToken = clientProps?.idToken
        }
        if (clientProps?.address) {
            this.address = clientProps?.address
        }
    }

    do<TRequest, TBody>(request: TRequest, url: string, method: Method, data: any, expectStatusCode: number): Promise<Response<TRequest, TBody>> {

        const headers: { [key: string]: string } = {}
        if (this.idToken) {
            headers["Authorization"] = `Bearer ${this.idToken}`
        }

        return axios.request<TBody>({
            method,
            url,
            data,
            headers
        }).then(value => {
            return clientResponse<TRequest, TBody>(value, request, expectStatusCode);
        }).catch((err) => {
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
        return this.do(request, `${this.address}/identityproviders`, "post", request.object, 200)
    }

    createOrganization(request: OrganizationCreateRequest): Promise<OrganizationCreateResponse> {
        return this.do(request, `${this.address}/organizations`, "post", request.object, 200)
    }

    createOAuth2Client(request: OAuth2ClientCreateRequest): Promise<OAuth2ClientCreateResponse> {
        return this.do(request, `${this.address}/clients`, "post", request.object, 200)
    }

    deleteOAuth2Client(request: OAuth2ClientDeleteRequest): Promise<OAuth2ClientDeleteResponse> {
        return this.do(request, `${this.address}/clients/${request.id}`, "delete", undefined, 204)
    }

    getIdentityProvider(request: IdentityProviderGetRequest): Promise<IdentityProviderGetResponse> {
        return this.do(request, `${this.address}/identityproviders/${request.id}`, "get", undefined, 200)
    }

    getOrganization(request: { id: string }): Promise<OrganizationGetResponse> {
        return this.do(request, `${this.address}/organizations/${request.id}`, "get", undefined, 200)
    }

    getOAuth2Client(request: OAuth2ClientGetRequest): Promise<OAuth2ClientGetResponse> {
        return this.do(request, `${this.address}/clients/${request.id}`, "get", undefined, 200)
    }

    listIdentityProviders(request: IdentityProviderListRequest): Promise<IdentityProviderListResponse> {
        return this.do(request, `${this.address}/identityproviders?organizationId=${request.organizationId}`, "get", undefined, 200)
    }

    listOAuth2Clients(request: OAuth2ClientListRequest): Promise<OAuth2ClientListResponse> {
        return this.do(request, `${this.address}/clients`, "get", undefined, 200)
    }

    listOrganizations(request: void): Promise<OrganizationListResponse> {
        return this.do(request, `${this.address}/organizations`, "get", undefined, 200)
    }

    updateIdentityProvider(request: IdentityProviderCreateRequest): Promise<IdentityProviderCreateResponse> {
        return this.do(request, `${this.address}/identityproviders/${request.object.id}`, "put", request.object, 200)
    }

    updateOAuth2Client(request: OAuth2ClientUpdateRequest): Promise<OAuth2ClientUpdateResponse> {
        return this.do(request, `${this.address}/clients/${request.object.id}`, "put", request.object, 200)
    }


}

export const defaultClient: Client = new client()

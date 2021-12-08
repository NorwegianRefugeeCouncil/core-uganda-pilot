import {
    OrganizationCreator,
    OrganizationLister,
    OrganizationGetter,
    IdentityProviderGetter,
    IdentityProviderLister,
    IdentityProviderCreator,
    IdentityProviderUpdater,
    OAuth2ClientGetter,
    OAuth2ClientLister,
    OAuth2ClientUpdater,
    OAuth2ClientCreator,
    OAuth2ClientDeleter,
    SessionGetter
} from "./restTypes";

export type ClientProps = {
    address?: string
}

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
        OAuth2ClientDeleter,
        SessionGetter {
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

export type Session = {
    active: boolean
    expiry: string
    expiredInSeconds: number
    subject: string
    username: string
}

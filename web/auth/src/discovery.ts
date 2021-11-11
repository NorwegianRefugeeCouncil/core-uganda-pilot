import axios from "axios";

export type DiscoveryDocument = {

    authorization_endpoint: string
    claims_parameter_supported: boolean
    claims_supported: string[]
    code_challenge_methods_supported: string[]
    end_session_endpoint: string
    grant_types_supported: string[]
    id_token_signing_alg_values_supported: string[]
    issuer: string
    jwks_uri: string
    response_modes_supported: string[]
    response_types_supported: string[]
    scopes_supported: string[]
    subject_types_supported: string[]
    token_endpoint_auth_methods_supported: string[]
    token_endpoint_auth_signing_alg_values_supported: string[]
    token_endpoint: string
    request_object_signing_alg_values_supported: string[]
    request_parameter_supported: boolean
    request_uri_parameter_supported: boolean
    require_request_uri_registration: boolean
    userinfo_endpoint: string
    introspection_endpoing: string
    introspection_endpoint_auth_methods_supported: string[]
    introspection_endpoint_auth_signing_alg_values_supported: string[]
    claim_types_supported: string[]
}

export type IssuerOrDiscovery = DiscoveryDocument | string

export async function resolveDiscoveryAsync(issuerOrDiscovery: IssuerOrDiscovery): Promise<DiscoveryDocument | null> {
    let issuer: string
    if (typeof issuerOrDiscovery === "string") {
        issuer = issuerOrDiscovery
    } else {
        issuer = issuerOrDiscovery.issuer
    }
    const metadataEndpoint = `${issuer}/.well-known/openid-configuration`;
    return axios.get<DiscoveryDocument>(metadataEndpoint)
        .then(value => value.data)
        .catch(err => {
            console.log("failed to get discovery document", err)
            return null
        })
}

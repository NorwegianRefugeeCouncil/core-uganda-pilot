import { CodedError } from "../utils/error";
import { TokenResponse } from "../utils/tokenrequest";
import React from "react";
import { AxiosInstance } from "axios";
export interface ResponseErrorConfig extends Record<string, any> {
    error: string;
    error_description?: string;
    error_uri?: string;
}
export interface AuthErrorConfig extends ResponseErrorConfig {
    state?: string;
}
export declare class ResponseError extends CodedError {
    /**
     * Used to assist the api-client developer in
     * understanding the error that occurred.
     */
    description?: string;
    /**
     * A URI identifying a human-readable web page with
     * information about the error, used to provide the api-client
     * developer with additional information about the error.
     */
    uri?: string;
    /**
     * Raw results of the error.
     */
    params: Record<string, string>;
    constructor(params: ResponseErrorConfig, errorCodeType: "auth" | "token");
}
/**
 * [Section 5.2](https://tools.ietf.org/html/rfc6749#section-5.2)
 */
export declare class AuthError extends ResponseError {
    /**
     * Required only if state is used in the initial request
     */
    state?: string;
    constructor(response: AuthErrorConfig);
}
/**
 * [Section 4.1.2.1](https://tools.ietf.org/html/rfc6749#section-4.1.2.1)
 */
export declare class TokenError extends ResponseError {
    constructor(response: ResponseErrorConfig);
}
export interface TokenResponseConfig {
    /**
     * The access token issued by the authorization server.
     *
     * [Section 4.2.2](https://tools.ietf.org/html/rfc6749#section-4.2.2)
     */
    accessToken: string;
    /**
     * The type of the token issued. Value is case insensitive.
     *
     * [Section 7.1](https://tools.ietf.org/html/rfc6749#section-7.1)
     */
    tokenType?: TokenType;
    /**
     * The lifetime in seconds of the access token.
     *
     * For example, the value `3600` denotes that the access token will
     * expire in one hour from the time the response was generated.
     *
     * If omitted, the authorization server should provide the
     * expiration time via other means or document the default value.
     *
     * [Section 4.2.2](https://tools.ietf.org/html/rfc6749#section-4.2.2)
     */
    expiresIn?: number;
    /**
     * The refresh token, which can be used to obtain new access tokens using the same authorization grant.
     *
     * [Section 5.1](https://tools.ietf.org/html/rfc6749#section-5.1)
     */
    refreshToken?: string;
    /**
     * The scope of the access token. Only required if it's different to the scope that was requested by the api-client.
     *
     * [Section 3.3](https://tools.ietf.org/html/rfc6749#section-3.3)
     */
    scope?: string;
    /**
     * Required if the "state" parameter was present in the api-client
     * authorization request.  The exact value received from the api-client.
     *
     * [Section 4.2.2](https://tools.ietf.org/html/rfc6749#section-4.2.2)
     */
    state?: string;
    /**
     * ID Token value associated with the authenticated session.
     *
     * [TokenResponse](https://openid.net/specs/openid-connect-core-1_0.html#TokenResponse)
     */
    idToken?: string;
    /**
     * Time in seconds when the token was received by the api-client.
     */
    issuedAt?: number;
}
export declare enum ResponseType {
    Code = "code",
    Token = "token",
    IdToken = "id_token"
}
export declare enum CodeChallengeMethod {
    S256 = "S256",
    Plain = "plain"
}
export declare enum Prompt {
    None = "none",
    Login = "login",
    Consent = "consent",
    SelectAccount = "select_account",
    Default = ""
}
export interface AuthRequestConfig {
    responseType?: ResponseType;
    clientId: string;
    redirectUri: string;
    scopes?: string[];
    codeChallengeMethod?: CodeChallengeMethod;
    codeChallenge?: string;
    prompt?: Prompt;
    state?: string;
    extraParams?: Record<string, string>;
    usePKCE?: boolean;
}
export declare type AuthSessionResult = {
    type: WebBrowserResultType.CANCEL | WebBrowserResultType.DISMISS | WebBrowserResultType.LOCKED;
} | {
    type: 'error' | 'success';
    errorCode: string | null;
    error?: AuthError | null;
    params: {
        [key: string]: string;
    };
    authentication: TokenResponse | null;
    url: string;
};
/**
 * Access token type
 *
 * [Section 7.1](https://tools.ietf.org/html/rfc6749#section-7.1)
 */
export declare type TokenType = 'bearer' | 'mac';
/**
 * A hint about the type of the token submitted for revocation.
 *
 * [Section 2.1](https://tools.ietf.org/html/rfc7009#section-2.1)
 */
export declare enum TokenTypeHint {
    /**
     * Access token.
     *
     * [Section 1.4](https://tools.ietf.org/html/rfc6749#section-1.4)
     */
    AccessToken = "access_token",
    /**
     * Refresh token.
     *
     * [Section 1.5](https://tools.ietf.org/html/rfc6749#section-1.5)
     */
    RefreshToken = "refresh_token"
}
export interface TokenRequestConfig {
    /**
     * A unique string representing the registration information provided by the api-client.
     * The api-client identifier is not a secret; it is exposed to the resource owner and shouldn't be used
     * alone for api-client authentication.
     *
     * The api-client identifier is unique to the authorization server.
     *
     * [Section 2.2](https://tools.ietf.org/html/rfc6749#section-2.2)
     */
    clientId: string;
    /**
     * ApiClient secret supplied by an auth provider.
     * There is no secure way to store this on the api-client.
     *
     * [Section 2.3.1](https://tools.ietf.org/html/rfc6749#section-2.3.1)
     */
    clientSecret?: string;
    /**
     * Extra query params that'll be added to the query string.
     */
    extraParams?: Record<string, string>;
    /**
     * List of strings to request access to.
     *
     * [Section 3.3](https://tools.ietf.org/html/rfc6749#section-3.3)
     */
    scopes?: string[];
}
export interface ServerTokenResponseConfig {
    access_token: string;
    token_type?: TokenType;
    expires_in?: number;
    refresh_token?: string;
    scope?: string;
    id_token?: string;
    issued_at?: number;
}
export interface AccessTokenRequestConfig extends TokenRequestConfig {
    /**
     * The authorization code received from the authorization server.
     */
    code: string;
    /**
     * If the `redirectUri` parameter was included in the `AuthRequest`, then it must be supplied here as well.
     *
     * [Section 3.1.2](https://tools.ietf.org/html/rfc6749#section-3.1.2)
     */
    redirectUri: string;
}
/**
 * Config used to request a token refresh, or code exchange.
 *
 * [Section 6](https://tools.ietf.org/html/rfc6749#section-6)
 */
export interface RefreshTokenRequestConfig extends TokenRequestConfig {
    /**
     * The refresh token issued to the api-client.
     */
    refreshToken?: string;
}
/**
 * Grant type values used in dynamic api-client registration and auth requests.
 *
 * [Appendix A.10](https://tools.ietf.org/html/rfc6749#appendix-A.10)
 */
export declare enum GrantType {
    /**
     * Used for exchanging an authorization code for one or more tokens.
     *
     * [Section 4.1.3](https://tools.ietf.org/html/rfc6749#section-4.1.3)
     */
    AuthorizationCode = "authorization_code",
    /**
     * Used when obtaining an access token.
     *
     * [Section 4.2](https://tools.ietf.org/html/rfc6749#section-4.2)
     */
    Implicit = "implicit",
    /**
     * Used when exchanging a refresh token for a new token.
     *
     * [Section 6](https://tools.ietf.org/html/rfc6749#section-6)
     */
    RefreshToken = "refresh_token",
    /**
     * Used for api-client credentials flow.
     *
     * [Section 4.4.2](https://tools.ietf.org/html/rfc6749#section-4.4.2)
     */
    ClientCredentials = "client_credentials"
}
/**
 * Options for the prompt window / web browser.
 * This can be used to configure how the web browser should look and behave.
 */
export declare type AuthRequestPromptOptions = {
    url?: string;
};
export declare type DiscoveryDocument = {
    authorization_endpoint: string;
    claims_parameter_supported: boolean;
    claims_supported: string[];
    code_challenge_methods_supported: string[];
    end_session_endpoint: string;
    grant_types_supported: string[];
    id_token_signing_alg_values_supported: string[];
    issuer: string;
    jwks_uri: string;
    response_modes_supported: string[];
    response_types_supported: string[];
    scopes_supported: string[];
    subject_types_supported: string[];
    token_endpoint_auth_methods_supported: string[];
    token_endpoint_auth_signing_alg_values_supported: string[];
    token_endpoint: string;
    request_object_signing_alg_values_supported: string[];
    request_parameter_supported: boolean;
    request_uri_parameter_supported: boolean;
    require_request_uri_registration: boolean;
    userinfo_endpoint: string;
    introspection_endpoing: string;
    introspection_endpoint_auth_methods_supported: string[];
    introspection_endpoint_auth_signing_alg_values_supported: string[];
    claim_types_supported: string[];
};
export declare type IssuerOrDiscovery = DiscoveryDocument | string;
export declare enum WebBrowserResultType {
    CANCEL = "cancel",
    DISMISS = "dismiss",
    OPENED = "opened",
    LOCKED = "locked",
    SUCCESS = "success"
}
export declare type WebBrowserAuthSessionResult = WebBrowserRedirectResult | WebBrowserResult;
export declare type WebBrowserRedirectResult = {
    type: WebBrowserResultType.SUCCESS;
    url: string;
};
export declare type WebBrowserResult = {
    type: WebBrowserResultType;
};
export declare type PromptMethod = (options?: AuthRequestPromptOptions) => Promise<AuthSessionResult>;
export declare type LoginComponentProps = {
    login: () => void;
};
export declare type AuthWrapperProps = {
    clientId: string;
    issuer: string;
    scopes?: string[];
    redirectUriSuffix?: string;
    customLoginComponent?: React.FC<LoginComponentProps>;
    handleLoginErr?: (err: any) => void;
    axiosInstance?: AxiosInstance;
};
//# sourceMappingURL=types.d.ts.map

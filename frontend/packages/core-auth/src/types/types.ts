import React from 'react';
import { AxiosInstance } from 'axios';

import { CodedError } from './error';
import { TokenResponse } from './response';

export interface ResponseErrorConfig extends Record<string, any> {
  error: string;
  error_description?: string;
  error_uri?: string;
}

export interface AuthErrorConfig extends ResponseErrorConfig {
  state?: string;
}

const errorCodeMessages: {
  auth: Record<string, string>;
  token: Record<string, string>;
} = {
  // https://tools.ietf.org/html/rfc6749#section-4.1.2.1
  // https://openid.net/specs/openid-connect-core-1_0.html#AuthError
  auth: {
    // OAuth 2.0
    invalid_request:
      'The request is missing a required parameter, includes an invalid parameter value, includes a parameter more than once, or is otherwise malformed.',
    unauthorized_client: 'The client is not authorized to request an authorization code using this method.',
    access_denied: 'The resource owner or authorization server denied the request.',
    unsupported_response_type: 'The authorization server does not support obtaining an authorization code using this method.',
    invalid_scope: 'The requested scope is invalid, unknown, or malformed.',
    server_error:
      'The authorization server encountered an unexpected condition that prevented it from fulfilling the request. (This error code is needed because a 500 Internal Server Error HTTP status code cannot be returned to the client via an HTTP redirect.)',
    temporarily_unavailable:
      'The authorization server is currently unable to handle the request due to a temporary overloading or maintenance of the server.  (This error code is needed because a 503 Service Unavailable HTTP status code cannot be returned to the client via an HTTP redirect.)',
    // Open ID Connect error codes
    interaction_required:
      'Auth server requires user interaction of some form to proceed. This error may be returned when the prompt parameter value in the auth request is none, but the auth request cannot be completed without displaying a user interface for user interaction.',
    login_required:
      'Auth server requires user authentication. This error may be returned when the prompt parameter value in the auth request is none, but the auth request cannot be completed without displaying a user interface for user authentication.',
    account_selection_required:
      'User is required to select a session at the auth server. The user may be authenticated at the auth server with different associated accounts, but the user did not select a session. This error may be returned when the prompt parameter value in the auth request is `none`, but the auth request cannot be completed without displaying a user interface to prompt for a session to use.',
    consent_required:
      'Auth server requires user consent. This error may be returned when the prompt parameter value in the auth request is none, but the auth request cannot be completed without displaying a user interface for user consent.',
    invalid_request_uri: 'The `request_uri` in the auth request returns an error or contains invalid data.',
    invalid_request_object: 'The request parameter contains an invalid request object.',
    request_not_supported:
      'The OP does not support use of the `request` parameter defined in Section 6. (https://openid.net/specs/openid-connect-core-1_0.html#JWTRequests)',
    request_uri_not_supported:
      'The OP does not support use of the `request_uri` parameter defined in Section 6. (https://openid.net/specs/openid-connect-core-1_0.html#JWTRequests)',
    registration_not_supported:
      'The OP does not support use of the `registration` parameter defined in Section 7.2.1. (https://openid.net/specs/openid-connect-core-1_0.html#RegistrationParameter)',
  },
  // https://tools.ietf.org/html/rfc6749#section-5.2
  token: {
    invalid_request:
      'The request is missing a required parameter, includes an unsupported parameter value (other than grant type), repeats a parameter, includes multiple credentials, utilizes more than one mechanism for authenticating the client, or is otherwise malformed.',
    invalid_client:
      'Client authentication failed (e.g., unknown client, no client authentication included, or unsupported authentication method).  The authorization server MAY return an HTTP 401 (Unauthorized) status code to indicate which HTTP authentication schemes are supported.  If the client attempted to authenticate via the "Authorization" request header field, the authorization server MUST respond with an HTTP 401 (Unauthorized) status code and include the "WWW-Authenticate" response header field matching the authentication scheme used by the client.',
    invalid_grant:
      'The provided authorization grant (e.g., authorization code, resource owner credentials) or refresh token is invalid, expired, revoked, does not match the redirection URI used in the authorization request, or was issued to another client.',
    unauthorized_client: 'The authenticated client is not authorized to use this authorization grant type.',
    unsupported_grant_type: 'The authorization grant type is not supported by the authorization server.',
  },
};

enum ResponseErrorCodeType {
  Auth = 'auth',
  Token = 'token',
}

export class ResponseError extends CodedError {
  /**
   * Used to assist the client developer in
   * understanding the error that occurred.
   */
  description?: string;

  /**
   * A URI identifying a human-readable web page with
   * information about the error, used to provide the client
   * developer with additional information about the error.
   */
  uri?: string;

  /**
   * Raw results of the error.
   */
  params: Record<string, string>;

  constructor(params: ResponseErrorConfig, errorCodeType: ResponseErrorCodeType) {
    const { error, error_description, error_uri } = params;
    const message = errorCodeMessages[errorCodeType][error];
    let errorMessage: string;
    if (message) {
      errorMessage = message + (error_description ? `\nMore info: ${error_description}` : '');
    } else if (error_description) {
      errorMessage = error_description;
    } else {
      errorMessage = 'An unknown error occurred';
    }
    super(error, errorMessage);
    this.description = error_description ?? message;
    this.uri = error_uri;
    this.params = params;
  }
}

/**
 * [Section 5.2](https://tools.ietf.org/html/rfc6749#section-5.2)
 */
export class AuthError extends ResponseError {
  /**
   * Required only if state is used in the initial request
   */
  state?: string;

  constructor(response: AuthErrorConfig) {
    super(response, ResponseErrorCodeType.Auth);
    this.state = response.state;
  }
}

/**
 * [Section 4.1.2.1](https://tools.ietf.org/html/rfc6749#section-4.1.2.1)
 */
export class TokenError extends ResponseError {
  constructor(response: ResponseErrorConfig) {
    super(response, ResponseErrorCodeType.Token);
  }
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
   * The scope of the access token. Only required if it's different to the scope that was requested by the client.
   *
   * [Section 3.3](https://tools.ietf.org/html/rfc6749#section-3.3)
   */
  scope?: string;
  /**
   * Required if the "state" parameter was present in the client
   * authorization request.  The exact value received from the client.
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
   * Time in seconds when the token was received by the client.
   */
  issuedAt?: number;
}

export enum ResponseType {
  Code = 'code',
  Token = 'token',
  IdToken = 'id_token',
}

export enum CodeChallengeMethod {
  S256 = 'S256',
  Plain = 'plain',
}

export enum Prompt {
  None = 'none',
  Login = 'login',
  Consent = 'consent',
  SelectAccount = 'select_account',
  Default = '',
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

export type AuthSessionResult =
  | {
      type: WebBrowserResultType.CANCEL | WebBrowserResultType.DISMISS | WebBrowserResultType.LOCKED;
    }
  | {
      type: 'error' | 'success';
      errorCode: string | null;
      error?: AuthError | null;
      params: { [key: string]: string };
      authentication: TokenResponse | null;
      url: string;
    };

/**
 * Access token type
 *
 * [Section 7.1](https://tools.ietf.org/html/rfc6749#section-7.1)
 */
export enum TokenType {
  Bearer = 'bearer',
  Mac = 'mac',
}

/**
 * A hint about the type of the token submitted for revocation.
 *
 * [Section 2.1](https://tools.ietf.org/html/rfc7009#section-2.1)
 */
export enum TokenTypeHint {
  /**
   * Access token.
   *
   * [Section 1.4](https://tools.ietf.org/html/rfc6749#section-1.4)
   */
  AccessToken = 'access_token',
  /**
   * Refresh token.
   *
   * [Section 1.5](https://tools.ietf.org/html/rfc6749#section-1.5)
   */
  RefreshToken = 'refresh_token',
}

export interface TokenRequestConfig {
  /**
   * A unique string representing the registration information provided by the client.
   * The client identifier is not a secret; it is exposed to the resource owner and shouldn't be used
   * alone for client authentication.
   *
   * The client identifier is unique to the authorization server.
   *
   * [Section 2.2](https://tools.ietf.org/html/rfc6749#section-2.2)
   */
  clientId: string;
  /**
   * Client secret supplied by an auth provider.
   * There is no secure way to store this on the client.
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
   * The refresh token issued to the client.
   */
  refreshToken?: string;
}

/**
 * Grant type values used in dynamic client registration and auth requests.
 *
 * [Appendix A.10](https://tools.ietf.org/html/rfc6749#appendix-A.10)
 */
export enum GrantType {
  /**
   * Used for exchanging an authorization code for one or more tokens.
   *
   * [Section 4.1.3](https://tools.ietf.org/html/rfc6749#section-4.1.3)
   */
  AuthorizationCode = 'authorization_code',
  /**
   * Used when obtaining an access token.
   *
   * [Section 4.2](https://tools.ietf.org/html/rfc6749#section-4.2)
   */
  Implicit = 'implicit',
  /**
   * Used when exchanging a refresh token for a new token.
   *
   * [Section 6](https://tools.ietf.org/html/rfc6749#section-6)
   */
  RefreshToken = 'refresh_token',
  /**
   * Used for client credentials flow.
   *
   * [Section 4.4.2](https://tools.ietf.org/html/rfc6749#section-4.4.2)
   */
  ClientCredentials = 'client_credentials',
}

/**
 * Options for the prompt window / web browser.
 * This can be used to configure how the web browser should look and behave.
 */
export type AuthRequestPromptOptions = {
  url?: string;
};

export type DiscoveryDocument = {
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

export type IssuerOrDiscovery = DiscoveryDocument | string;

export enum WebBrowserResultType {
  CANCEL = 'cancel',
  DISMISS = 'dismiss',
  OPENED = 'opened',
  LOCKED = 'locked',
  SUCCESS = 'success',
}

export type WebBrowserAuthSessionResult = WebBrowserRedirectResult | WebBrowserResult;

export type WebBrowserRedirectResult = {
  type: WebBrowserResultType.SUCCESS;
  url: string;
};

export type WebBrowserResult = {
  type: WebBrowserResultType;
};

export type PromptMethod = (options?: AuthRequestPromptOptions) => Promise<AuthSessionResult>;

export type LoginComponentProps = {
  login: () => void;
};

export type AuthWrapperProps = {
  clientId: string;
  issuer: string;
  scopes: string[];
  redirectUri: string;
  customLoginComponent?: React.FC<LoginComponentProps>;
  handleLoginErr?: (err: any) => void;
  axiosInstance?: AxiosInstance;
  injectToken?: 'access_token' | 'id_token';
};

export type listenerMapEntry = {
  listener: (event: MessageEvent) => void;
  interval: ReturnType<typeof setTimeout>;
};

/* eslint-disable max-classes-per-file */
import axios, { AxiosRequestHeaders } from 'axios';

import {
  AccessTokenRequestConfig,
  DiscoveryDocument,
  GrantType,
  RefreshTokenRequestConfig,
  ResponseErrorConfig,
  ServerTokenResponseConfig,
  TokenError,
  TokenRequestConfig,
} from './types';
import { TokenResponse } from './response';

class Request<T, B> {
  constructor(protected request: T) {}

  async performAsync(discovery: DiscoveryDocument): Promise<B> {
    throw new Error('performAsync must be extended');
  }

  getRequestConfig(): T {
    throw new Error('getRequestConfig must be extended');
  }

  getQueryBody(): URLSearchParams {
    throw new Error('getQueryBody must be extended');
  }
}

class TokenRequest<T extends TokenRequestConfig> extends Request<
  T,
  TokenResponse
> {
  readonly clientId: string;

  readonly scopes?: string[];

  readonly extraParams?: Record<string, string>;

  constructor(request: T, public grantType: GrantType) {
    super(request);
    this.clientId = request.clientId;
    this.extraParams = request.extraParams;
    this.scopes = request.scopes;
  }

  getHeaders(): AxiosRequestHeaders {
    const header: Record<string, string> = {};
    header['Content-Type'] = 'application/x-www-form-urlencoded';
    return header;
  }

  async performAsync(discovery: Pick<DiscoveryDocument, 'token_endpoint'>) {
    const response = await axios
      .post<ServerTokenResponseConfig>(
        discovery.token_endpoint,
        this.getQueryBody().toString(),
        {
          headers: this.getHeaders(),
        },
      )
      .catch((err) => {
        if (axios.isAxiosError(err)) {
          if (
            err.response &&
            err.response.data &&
            'error' in (err.response.data as ResponseErrorConfig)
          ) {
            throw new TokenError(err.response.data as ResponseErrorConfig);
          }
        }
        throw err;
      });

    return new TokenResponse({
      accessToken: response.data.access_token,
      tokenType: response.data.token_type,
      expiresIn: response.data.expires_in,
      refreshToken: response.data.refresh_token,
      scope: response.data.scope,
      idToken: response.data.id_token,
      issuedAt: response.data.issued_at,
    });
  }

  getQueryBody(): URLSearchParams {
    const queryBody = new URLSearchParams();
    queryBody.set('grant_type', this.grantType);

    if (this.scopes) {
      queryBody.set('scope', this.scopes.join(' '));
    }

    queryBody.set('client_id', this.clientId);

    if (this.extraParams) {
      for (const extra in this.extraParams) {
        if (extra in this.extraParams && !(extra in queryBody)) {
          queryBody.set(extra, this.extraParams[extra]);
        }
      }
    }
    return queryBody;
  }
}

/**
 * Access token request. Exchange an authorization code for a user access token.
 *
 * [Section 4.1.3](https://tools.ietf.org/html/rfc6749#section-4.1.3)
 */
export class AccessTokenRequest
  extends TokenRequest<AccessTokenRequestConfig>
  implements AccessTokenRequestConfig
{
  readonly code: string;

  readonly redirectUri: string;

  constructor(options: AccessTokenRequestConfig) {
    super(options, GrantType.AuthorizationCode);
    this.code = options.code;
    this.redirectUri = options.redirectUri;
  }

  getQueryBody(): URLSearchParams {
    const queryBody: URLSearchParams = super.getQueryBody();

    if (this.redirectUri) {
      queryBody.set('redirect_uri', this.redirectUri);
    }

    if (this.code) {
      queryBody.set('code', this.code);
    }

    return queryBody;
  }

  getRequestConfig() {
    return {
      clientId: this.clientId,
      grantType: this.grantType,
      code: this.code,
      redirectUri: this.redirectUri,
      extraParams: this.extraParams,
      scopes: this.scopes,
    };
  }
}

/**
 * Refresh request.
 *
 * [Section 6](https://tools.ietf.org/html/rfc6749#section-6)
 */
export class RefreshTokenRequest
  extends TokenRequest<RefreshTokenRequestConfig>
  implements RefreshTokenRequestConfig
{
  readonly refreshToken?: string;

  constructor(options: RefreshTokenRequestConfig) {
    super(options, GrantType.RefreshToken);
    this.refreshToken = options.refreshToken;
  }

  getQueryBody() {
    const queryBody = super.getQueryBody();
    if (this.refreshToken) {
      queryBody.set('refresh_token', this.refreshToken);
    }
    return queryBody;
  }

  getRequestConfig() {
    return {
      clientId: this.clientId,
      grantType: this.grantType,
      refreshToken: this.refreshToken,
      extraParams: this.extraParams,
      scopes: this.scopes,
    };
  }
}

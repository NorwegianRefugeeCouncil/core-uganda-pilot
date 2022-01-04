import getCurrentTimeInSeconds from '../utils/getCurrentTimeInSeconds';

import {
  DiscoveryDocument,
  TokenRequestConfig,
  TokenResponseConfig,
  TokenType,
} from './types';
import { RefreshTokenRequest } from './request';

export class TokenResponse implements TokenResponseConfig {
  accessToken: string;

  tokenType: TokenType;

  expiresIn?: number;

  refreshToken?: string;

  scope?: string;

  state?: string;

  idToken?: string;

  issuedAt: number;

  public constructor(response: TokenResponseConfig) {
    this.accessToken = response.accessToken;
    this.tokenType = response.tokenType ?? TokenType.Bearer;
    this.expiresIn = response.expiresIn;
    this.refreshToken = response.refreshToken;
    this.scope = response.scope;
    this.state = response.state;
    this.idToken = response.idToken;
    this.issuedAt = response.issuedAt ?? getCurrentTimeInSeconds();
  }

  static isTokenFresh(
    token: Pick<TokenResponse, 'expiresIn' | 'issuedAt'>,
    marginPercentage = 0.1,
  ): boolean {
    if (!token) return false;

    if (token.expiresIn) {
      const secondsMargin = token.expiresIn * marginPercentage * -1;
      const now = getCurrentTimeInSeconds();
      return now < token.issuedAt + token.expiresIn + secondsMargin;
    }

    return true;
  }

  static fromQueryParams(params: Record<string, string>): TokenResponse {
    return new TokenResponse({
      accessToken: params.access_token,
      refreshToken: params.refresh_token,
      scope: params.scope,
      state: params.state,
      idToken: params.id_token,
      tokenType: params.token_type as TokenType,
      expiresIn: params.expires_in ? parseInt(params.expires_in) : undefined,
      issuedAt: params.expires_in ? parseInt(params.issuedAt) : undefined,
    });
  }

  private applyResponseConfig(response: TokenResponseConfig) {
    this.accessToken = response.accessToken ?? this.accessToken;
    this.tokenType = response.tokenType ?? this.tokenType ?? TokenType.Bearer;
    this.expiresIn = response.expiresIn ?? this.expiresIn;
    this.refreshToken = response.refreshToken ?? this.refreshToken;
    this.scope = response.scope ?? this.scope;
    this.state = response.state ?? this.state;
    this.idToken = response.idToken ?? this.idToken;
    this.issuedAt =
      response.issuedAt ?? this.issuedAt ?? getCurrentTimeInSeconds();
  }

  getRequestConfig(): TokenResponseConfig {
    return {
      accessToken: this.accessToken,
      idToken: this.idToken,
      refreshToken: this.refreshToken,
      scope: this.scope,
      state: this.state,
      tokenType: this.tokenType,
      issuedAt: this.issuedAt,
      expiresIn: this.expiresIn,
    };
  }

  shouldRefresh(): boolean {
    return !(TokenResponse.isTokenFresh(this) || !this.refreshToken);
  }

  getExpiryMs(): number {
    const tenPercentMargin = 0.9;
    return 1000 * (this.issuedAt + (this.expiresIn ?? 0) * tenPercentMargin);
  }

  async refreshAsync(
    config: Omit<TokenRequestConfig, 'grantType' | 'refreshToken'>,
    discovery: Pick<DiscoveryDocument, 'token_endpoint'>,
  ): Promise<this> {
    const request = new RefreshTokenRequest({
      ...config,
      refreshToken: this.refreshToken,
    });
    let response: TokenResponse;
    try {
      response = await request.performAsync(discovery);
    } catch (e) {
      throw new Error(`Error encountered while performing token refresh: ${e}`);
    }

    response.refreshToken = response.refreshToken ?? this.refreshToken;
    const json = response.getRequestConfig();
    this.applyResponseConfig(json);

    return this;
  }
}

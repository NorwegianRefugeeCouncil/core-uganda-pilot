import getCurrentTimeInSeconds from '../utils/getCurrentTimeInSeconds';

import {
  DiscoveryDocument,
  RefreshTokenRequestConfig,
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
      expiresIn: params.expires_in
        ? parseInt(params.expires_in, 10)
        : undefined,
      issuedAt: params.expires_in ? parseInt(params.issuedAt, 10) : undefined,
    });
  }

  static createTokenResponse = (
    trc: TokenResponseConfig,
  ): TokenResponse | undefined => {
    if (!trc) return undefined;
    try {
      return new TokenResponse(trc);
    } catch (e) {
      return undefined;
    }
  };

  static async refreshAsync(
    config: Omit<RefreshTokenRequestConfig, 'grantType'>,
    discovery: Pick<DiscoveryDocument, 'token_endpoint'>,
  ): Promise<TokenResponse> {
    const request = new RefreshTokenRequest(config);
    let response: TokenResponse;
    try {
      response = await request.performAsync(discovery);
    } catch (e) {
      throw new Error(`Error encountered while performing token refresh: ${e}`);
    }

    response.refreshToken = response.refreshToken ?? config.refreshToken;
    const json = response.getRequestConfig();

    return new TokenResponse(json);
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
}

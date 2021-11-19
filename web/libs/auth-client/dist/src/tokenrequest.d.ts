import { AccessTokenRequestConfig, DiscoveryDocument, GrantType, RefreshTokenRequestConfig, TokenRequestConfig, TokenResponseConfig, TokenType } from "./types/types";
import { AxiosRequestHeaders } from "axios";
export declare function getCurrentTimeInSeconds(): number;
export declare class TokenResponse implements TokenResponseConfig {
    static isTokenFresh(token: Pick<TokenResponse, "expiresIn" | "issuedAt">, secondsMargin?: number): boolean;
    static fromQueryParams(params: Record<string, string>): TokenResponse;
    accessToken: string;
    tokenType: TokenType;
    expiresIn?: number;
    refreshToken?: string;
    scope?: string;
    state?: string;
    idToken?: string;
    issuedAt: number;
    constructor(response: TokenResponseConfig);
    private applyResponseConfig;
    getRequestConfig(): TokenResponseConfig;
    shouldRefresh(): boolean;
    refreshAsync(config: Omit<TokenRequestConfig, "grantType" | "refreshToken">, discovery: Pick<DiscoveryDocument, "token_endpoint">): Promise<this>;
}
declare class Request<T, B> {
    protected request: T;
    constructor(request: T);
    performAsync(discovery: DiscoveryDocument): Promise<B>;
    getRequestConfig(): T;
    getQueryBody(): URLSearchParams;
}
declare class TokenRequest<T extends TokenRequestConfig> extends Request<T, TokenResponse> {
    grantType: GrantType;
    readonly clientId: string;
    readonly scopes?: string[];
    readonly extraParams?: Record<string, string>;
    constructor(request: T, grantType: GrantType);
    getHeaders(): AxiosRequestHeaders;
    performAsync(discovery: Pick<DiscoveryDocument, 'token_endpoint'>): Promise<TokenResponse>;
    getQueryBody(): URLSearchParams;
}
/**
 * Access token request. Exchange an authorization code for a user access token.
 *
 * [Section 4.1.3](https://tools.ietf.org/html/rfc6749#section-4.1.3)
 */
export declare class AccessTokenRequest extends TokenRequest<AccessTokenRequestConfig> implements AccessTokenRequestConfig {
    readonly code: string;
    readonly redirectUri: string;
    constructor(options: AccessTokenRequestConfig);
    getQueryBody(): URLSearchParams;
    getRequestConfig(): {
        clientId: string;
        grantType: GrantType;
        code: string;
        redirectUri: string;
        extraParams: Record<string, string>;
        scopes: string[];
    };
}
/**
 * Refresh request.
 *
 * [Section 6](https://tools.ietf.org/html/rfc6749#section-6)
 */
export declare class RefreshTokenRequest extends TokenRequest<RefreshTokenRequestConfig> implements RefreshTokenRequestConfig {
    readonly refreshToken?: string;
    constructor(options: RefreshTokenRequestConfig);
    getQueryBody(): URLSearchParams;
    getRequestConfig(): {
        clientId: string;
        grantType: GrantType;
        refreshToken: string;
        extraParams: Record<string, string>;
        scopes: string[];
    };
}
export declare function exchangeCodeAsync(config: AccessTokenRequestConfig, discovery: Pick<DiscoveryDocument, 'token_endpoint'>): Promise<TokenResponse>;
export {};

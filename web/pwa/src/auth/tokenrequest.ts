import {
    AccessTokenRequestConfig,
    GrantType,
    RefreshTokenRequestConfig,
    ResponseErrorConfig,
    ServerTokenResponseConfig,
    TokenError,
    TokenRequestConfig,
    TokenResponseConfig,
    TokenType
} from "./types";
import {DiscoveryDocument} from "./discovery";
import axios, {AxiosRequestHeaders, AxiosResponse} from "axios"
import qs from "qs";


export function getCurrentTimeInSeconds(): number {
    return Math.floor(Date.now() / 1000);
}

export class TokenResponse implements TokenResponseConfig {

    static isTokenFresh(token: Pick<TokenResponse, "expiresIn" | "issuedAt">, secondsMargin: number = 60 * 10 - 1): boolean {
        if (!token) {
            return false
        }
        if (token.expiresIn) {
            const now = getCurrentTimeInSeconds()
            return now < token.issuedAt + token.expiresIn + secondsMargin
        }
        return true
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

    accessToken: string;
    tokenType: TokenType;
    expiresIn?: number;
    refreshToken?: string;
    scope?: string;
    state?: string;
    idToken?: string;
    issuedAt: number;

    public constructor(response: TokenResponseConfig) {
        this.accessToken = response.accessToken
        this.tokenType = response.tokenType ?? "bearer"
        this.expiresIn = response.expiresIn
        this.refreshToken = response.refreshToken
        this.scope = response.scope
        this.state = response.state
        this.idToken = response.idToken
        this.issuedAt = response.issuedAt ?? getCurrentTimeInSeconds()
    }

    private applyResponseConfig(response: TokenResponseConfig) {
        this.accessToken = response.accessToken ?? this.accessToken;
        this.tokenType = response.tokenType ?? this.tokenType ?? 'bearer';
        this.expiresIn = response.expiresIn ?? this.expiresIn;
        this.refreshToken = response.refreshToken ?? this.refreshToken;
        this.scope = response.scope ?? this.scope;
        this.state = response.state ?? this.state;
        this.idToken = response.idToken ?? this.idToken;
        this.issuedAt = response.issuedAt ?? this.issuedAt ?? getCurrentTimeInSeconds();
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

    async refreshAsync(
        config: Omit<TokenRequestConfig, "grantType" | "refreshToken">,
        discovery: Pick<DiscoveryDocument, "token_endpoint">) {
        const request = new RefreshTokenRequest({...config, refreshToken: this.refreshToken})
        const response = await request.performAsync(discovery)
        response.refreshToken = response.refreshToken ?? this.refreshToken
        const json = response.getRequestConfig()
        this.applyResponseConfig(json)
        return this
    }

}

class Request<T, B> {
    constructor(protected request: T) {
    }

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

class TokenRequest<T extends TokenRequestConfig> extends Request<T, TokenResponse> {
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
        const header: Record<string, string> = {}
        header["Content-Type"] = "application/x-www-form-urlencoded"
        return header;
    }

    async performAsync(discovery: Pick<DiscoveryDocument, 'token_endpoint'>) {

        const response = await axios({
            method: "POST",
            headers: this.getHeaders(),
            data: this.getQueryBody().toString(),
            url: discovery.token_endpoint,
        }).catch(err => {
            if (axios.isAxiosError(err)) {
                if (err.response && err.response.data && "error" in (err.response.data as ResponseErrorConfig)) {
                    throw new TokenError(err.response.data as ResponseErrorConfig);
                }
            }
            throw err
        }) as AxiosResponse<ServerTokenResponseConfig>;

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

        const queryBody = new URLSearchParams()
        queryBody.set("grant_type", this.grantType)

        if (this.scopes) {
            queryBody.set("scope", this.scopes.join(" "))
        }

        queryBody.set("client_id", this.clientId)

        if (this.extraParams) {
            for (const extra in this.extraParams) {
                if (extra in this.extraParams && !(extra in queryBody)) {
                    queryBody.set(extra, this.extraParams[extra])
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
    implements AccessTokenRequestConfig {
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
            queryBody.set("redirect_uri", this.redirectUri)
        }

        if (this.code) {
            queryBody.set("code", this.code)
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
    implements RefreshTokenRequestConfig {
    readonly refreshToken?: string;

    constructor(options: RefreshTokenRequestConfig) {
        super(options, GrantType.RefreshToken);
        this.refreshToken = options.refreshToken;
    }

    getQueryBody() {
        const queryBody = super.getQueryBody();
        if (this.refreshToken) {
            queryBody.set("refresh_token", this.refreshToken)
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

export function exchangeCodeAsync(
    config: AccessTokenRequestConfig,
    discovery: Pick<DiscoveryDocument, 'token_endpoint'>
): Promise<TokenResponse> {
    const request = new AccessTokenRequest(config);
    console.log("ACCESS TOKEN EXCHANGE CONFIG", config)
    return request.performAsync(discovery);
}


export function refreshAsync(
    config: RefreshTokenRequestConfig,
    discovery: Pick<DiscoveryDocument, 'token_endpoint'>
): Promise<TokenResponse> {
    const request = new RefreshTokenRequest(config);
    return request.performAsync(discovery);
}

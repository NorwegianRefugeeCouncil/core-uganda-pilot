"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.exchangeCodeAsync = exports.RefreshTokenRequest = exports.AccessTokenRequest = exports.TokenResponse = exports.getCurrentTimeInSeconds = void 0;
const types_1 = require("./types/types");
const axios_1 = __importDefault(require("axios"));
function getCurrentTimeInSeconds() {
    return Math.floor(Date.now() / 1000);
}
exports.getCurrentTimeInSeconds = getCurrentTimeInSeconds;
class TokenResponse {
    constructor(response) {
        var _a, _b;
        this.accessToken = response.accessToken;
        this.tokenType = (_a = response.tokenType) !== null && _a !== void 0 ? _a : "bearer";
        this.expiresIn = response.expiresIn;
        this.refreshToken = response.refreshToken;
        this.scope = response.scope;
        this.state = response.state;
        this.idToken = response.idToken;
        this.issuedAt = (_b = response.issuedAt) !== null && _b !== void 0 ? _b : getCurrentTimeInSeconds();
    }
    static isTokenFresh(token, secondsMargin = 60 * 10 - 1) {
        if (!token) {
            return false;
        }
        if (token.expiresIn) {
            const now = getCurrentTimeInSeconds();
            return now < token.issuedAt + token.expiresIn + secondsMargin;
        }
        return true;
    }
    static fromQueryParams(params) {
        return new TokenResponse({
            accessToken: params.access_token,
            refreshToken: params.refresh_token,
            scope: params.scope,
            state: params.state,
            idToken: params.id_token,
            tokenType: params.token_type,
            expiresIn: params.expires_in ? parseInt(params.expires_in) : undefined,
            issuedAt: params.expires_in ? parseInt(params.issuedAt) : undefined,
        });
    }
    applyResponseConfig(response) {
        var _a, _b, _c, _d, _e, _f, _g, _h, _j, _k;
        this.accessToken = (_a = response.accessToken) !== null && _a !== void 0 ? _a : this.accessToken;
        this.tokenType = (_c = (_b = response.tokenType) !== null && _b !== void 0 ? _b : this.tokenType) !== null && _c !== void 0 ? _c : 'bearer';
        this.expiresIn = (_d = response.expiresIn) !== null && _d !== void 0 ? _d : this.expiresIn;
        this.refreshToken = (_e = response.refreshToken) !== null && _e !== void 0 ? _e : this.refreshToken;
        this.scope = (_f = response.scope) !== null && _f !== void 0 ? _f : this.scope;
        this.state = (_g = response.state) !== null && _g !== void 0 ? _g : this.state;
        this.idToken = (_h = response.idToken) !== null && _h !== void 0 ? _h : this.idToken;
        this.issuedAt = (_k = (_j = response.issuedAt) !== null && _j !== void 0 ? _j : this.issuedAt) !== null && _k !== void 0 ? _k : getCurrentTimeInSeconds();
    }
    getRequestConfig() {
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
    shouldRefresh() {
        return !(TokenResponse.isTokenFresh(this) || !this.refreshToken);
    }
    refreshAsync(config, discovery) {
        var _a;
        return __awaiter(this, void 0, void 0, function* () {
            const request = new RefreshTokenRequest(Object.assign(Object.assign({}, config), { refreshToken: this.refreshToken }));
            const response = yield request.performAsync(discovery);
            response.refreshToken = (_a = response.refreshToken) !== null && _a !== void 0 ? _a : this.refreshToken;
            const json = response.getRequestConfig();
            this.applyResponseConfig(json);
            return this;
        });
    }
}
exports.TokenResponse = TokenResponse;
class Request {
    constructor(request) {
        this.request = request;
    }
    performAsync(discovery) {
        return __awaiter(this, void 0, void 0, function* () {
            throw new Error('performAsync must be extended');
        });
    }
    getRequestConfig() {
        throw new Error('getRequestConfig must be extended');
    }
    getQueryBody() {
        throw new Error('getQueryBody must be extended');
    }
}
class TokenRequest extends Request {
    constructor(request, grantType) {
        super(request);
        this.grantType = grantType;
        this.clientId = request.clientId;
        this.extraParams = request.extraParams;
        this.scopes = request.scopes;
    }
    getHeaders() {
        const header = {};
        header["Content-Type"] = "application/x-www-form-urlencoded";
        return header;
    }
    performAsync(discovery) {
        return __awaiter(this, void 0, void 0, function* () {
            const response = yield (0, axios_1.default)({
                method: "POST",
                headers: this.getHeaders(),
                data: this.getQueryBody().toString(),
                url: discovery.token_endpoint,
            }).catch(err => {
                if (axios_1.default.isAxiosError(err)) {
                    if (err.response && err.response.data && "error" in err.response.data) {
                        throw new types_1.TokenError(err.response.data);
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
        });
    }
    getQueryBody() {
        const queryBody = new URLSearchParams();
        queryBody.set("grant_type", this.grantType);
        if (this.scopes) {
            queryBody.set("scope", this.scopes.join(" "));
        }
        queryBody.set("client_id", this.clientId);
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
class AccessTokenRequest extends TokenRequest {
    constructor(options) {
        super(options, types_1.GrantType.AuthorizationCode);
        this.code = options.code;
        this.redirectUri = options.redirectUri;
    }
    getQueryBody() {
        const queryBody = super.getQueryBody();
        if (this.redirectUri) {
            queryBody.set("redirect_uri", this.redirectUri);
        }
        if (this.code) {
            queryBody.set("code", this.code);
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
exports.AccessTokenRequest = AccessTokenRequest;
/**
 * Refresh request.
 *
 * [Section 6](https://tools.ietf.org/html/rfc6749#section-6)
 */
class RefreshTokenRequest extends TokenRequest {
    constructor(options) {
        super(options, types_1.GrantType.RefreshToken);
        this.refreshToken = options.refreshToken;
    }
    getQueryBody() {
        const queryBody = super.getQueryBody();
        if (this.refreshToken) {
            queryBody.set("refresh_token", this.refreshToken);
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
exports.RefreshTokenRequest = RefreshTokenRequest;
function exchangeCodeAsync(config, discovery) {
    const request = new AccessTokenRequest(config);
    console.log("ACCESS TOKEN EXCHANGE CONFIG", config);
    return request.performAsync(discovery);
}
exports.exchangeCodeAsync = exchangeCodeAsync;

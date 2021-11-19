"use strict";
var __extends = (this && this.__extends) || (function () {
    var extendStatics = function (d, b) {
        extendStatics = Object.setPrototypeOf ||
            ({ __proto__: [] } instanceof Array && function (d, b) { d.__proto__ = b; }) ||
            function (d, b) { for (var p in b) if (Object.prototype.hasOwnProperty.call(b, p)) d[p] = b[p]; };
        return extendStatics(d, b);
    };
    return function (d, b) {
        if (typeof b !== "function" && b !== null)
            throw new TypeError("Class extends value " + String(b) + " is not a constructor or null");
        extendStatics(d, b);
        function __() { this.constructor = d; }
        d.prototype = b === null ? Object.create(b) : (__.prototype = b.prototype, new __());
    };
})();
var __assign = (this && this.__assign) || function () {
    __assign = Object.assign || function(t) {
        for (var s, i = 1, n = arguments.length; i < n; i++) {
            s = arguments[i];
            for (var p in s) if (Object.prototype.hasOwnProperty.call(s, p))
                t[p] = s[p];
        }
        return t;
    };
    return __assign.apply(this, arguments);
};
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
var __generator = (this && this.__generator) || function (thisArg, body) {
    var _ = { label: 0, sent: function() { if (t[0] & 1) throw t[1]; return t[1]; }, trys: [], ops: [] }, f, y, t, g;
    return g = { next: verb(0), "throw": verb(1), "return": verb(2) }, typeof Symbol === "function" && (g[Symbol.iterator] = function() { return this; }), g;
    function verb(n) { return function (v) { return step([n, v]); }; }
    function step(op) {
        if (f) throw new TypeError("Generator is already executing.");
        while (_) try {
            if (f = 1, y && (t = op[0] & 2 ? y["return"] : op[0] ? y["throw"] || ((t = y["return"]) && t.call(y), 0) : y.next) && !(t = t.call(y, op[1])).done) return t;
            if (y = 0, t) op = [op[0] & 2, t.value];
            switch (op[0]) {
                case 0: case 1: t = op; break;
                case 4: _.label++; return { value: op[1], done: false };
                case 5: _.label++; y = op[1]; op = [0]; continue;
                case 7: op = _.ops.pop(); _.trys.pop(); continue;
                default:
                    if (!(t = _.trys, t = t.length > 0 && t[t.length - 1]) && (op[0] === 6 || op[0] === 2)) { _ = 0; continue; }
                    if (op[0] === 3 && (!t || (op[1] > t[0] && op[1] < t[3]))) { _.label = op[1]; break; }
                    if (op[0] === 6 && _.label < t[1]) { _.label = t[1]; t = op; break; }
                    if (t && _.label < t[2]) { _.label = t[2]; _.ops.push(op); break; }
                    if (t[2]) _.ops.pop();
                    _.trys.pop(); continue;
            }
            op = body.call(thisArg, _);
        } catch (e) { op = [6, e]; y = 0; } finally { f = t = 0; }
        if (op[0] & 5) throw op[1]; return { value: op[0] ? op[1] : void 0, done: true };
    }
};
exports.__esModule = true;
exports.exchangeCodeAsync = exports.RefreshTokenRequest = exports.AccessTokenRequest = exports.TokenResponse = exports.getCurrentTimeInSeconds = void 0;
var types_1 = require("./types/types");
var axios_1 = require("axios");
function getCurrentTimeInSeconds() {
    return Math.floor(Date.now() / 1000);
}
exports.getCurrentTimeInSeconds = getCurrentTimeInSeconds;
var TokenResponse = /** @class */ (function () {
    function TokenResponse(response) {
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
    TokenResponse.isTokenFresh = function (token, secondsMargin) {
        if (secondsMargin === void 0) { secondsMargin = 60 * 10 - 1; }
        if (!token) {
            return false;
        }
        if (token.expiresIn) {
            var now = getCurrentTimeInSeconds();
            return now < token.issuedAt + token.expiresIn + secondsMargin;
        }
        return true;
    };
    TokenResponse.fromQueryParams = function (params) {
        return new TokenResponse({
            accessToken: params.access_token,
            refreshToken: params.refresh_token,
            scope: params.scope,
            state: params.state,
            idToken: params.id_token,
            tokenType: params.token_type,
            expiresIn: params.expires_in ? parseInt(params.expires_in) : undefined,
            issuedAt: params.expires_in ? parseInt(params.issuedAt) : undefined
        });
    };
    TokenResponse.prototype.applyResponseConfig = function (response) {
        var _a, _b, _c, _d, _e, _f, _g, _h, _j, _k;
        this.accessToken = (_a = response.accessToken) !== null && _a !== void 0 ? _a : this.accessToken;
        this.tokenType = (_c = (_b = response.tokenType) !== null && _b !== void 0 ? _b : this.tokenType) !== null && _c !== void 0 ? _c : 'bearer';
        this.expiresIn = (_d = response.expiresIn) !== null && _d !== void 0 ? _d : this.expiresIn;
        this.refreshToken = (_e = response.refreshToken) !== null && _e !== void 0 ? _e : this.refreshToken;
        this.scope = (_f = response.scope) !== null && _f !== void 0 ? _f : this.scope;
        this.state = (_g = response.state) !== null && _g !== void 0 ? _g : this.state;
        this.idToken = (_h = response.idToken) !== null && _h !== void 0 ? _h : this.idToken;
        this.issuedAt = (_k = (_j = response.issuedAt) !== null && _j !== void 0 ? _j : this.issuedAt) !== null && _k !== void 0 ? _k : getCurrentTimeInSeconds();
    };
    TokenResponse.prototype.getRequestConfig = function () {
        return {
            accessToken: this.accessToken,
            idToken: this.idToken,
            refreshToken: this.refreshToken,
            scope: this.scope,
            state: this.state,
            tokenType: this.tokenType,
            issuedAt: this.issuedAt,
            expiresIn: this.expiresIn
        };
    };
    TokenResponse.prototype.shouldRefresh = function () {
        return !(TokenResponse.isTokenFresh(this) || !this.refreshToken);
    };
    TokenResponse.prototype.refreshAsync = function (config, discovery) {
        var _a;
        return __awaiter(this, void 0, void 0, function () {
            var request, response, json;
            return __generator(this, function (_b) {
                switch (_b.label) {
                    case 0:
                        request = new RefreshTokenRequest(__assign(__assign({}, config), { refreshToken: this.refreshToken }));
                        return [4 /*yield*/, request.performAsync(discovery)];
                    case 1:
                        response = _b.sent();
                        response.refreshToken = (_a = response.refreshToken) !== null && _a !== void 0 ? _a : this.refreshToken;
                        json = response.getRequestConfig();
                        this.applyResponseConfig(json);
                        return [2 /*return*/, this];
                }
            });
        });
    };
    return TokenResponse;
}());
exports.TokenResponse = TokenResponse;
var Request = /** @class */ (function () {
    function Request(request) {
        this.request = request;
    }
    Request.prototype.performAsync = function (discovery) {
        return __awaiter(this, void 0, void 0, function () {
            return __generator(this, function (_a) {
                throw new Error('performAsync must be extended');
            });
        });
    };
    Request.prototype.getRequestConfig = function () {
        throw new Error('getRequestConfig must be extended');
    };
    Request.prototype.getQueryBody = function () {
        throw new Error('getQueryBody must be extended');
    };
    return Request;
}());
var TokenRequest = /** @class */ (function (_super) {
    __extends(TokenRequest, _super);
    function TokenRequest(request, grantType) {
        var _this = _super.call(this, request) || this;
        _this.grantType = grantType;
        _this.clientId = request.clientId;
        _this.extraParams = request.extraParams;
        _this.scopes = request.scopes;
        return _this;
    }
    TokenRequest.prototype.getHeaders = function () {
        var header = {};
        header["Content-Type"] = "application/x-www-form-urlencoded";
        return header;
    };
    TokenRequest.prototype.performAsync = function (discovery) {
        return __awaiter(this, void 0, void 0, function () {
            var response;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0: return [4 /*yield*/, (0, axios_1["default"])({
                            method: "POST",
                            headers: this.getHeaders(),
                            data: this.getQueryBody().toString(),
                            url: discovery.token_endpoint
                        })["catch"](function (err) {
                            if (axios_1["default"].isAxiosError(err)) {
                                if (err.response && err.response.data && "error" in err.response.data) {
                                    throw new types_1.TokenError(err.response.data);
                                }
                            }
                            throw err;
                        })];
                    case 1:
                        response = _a.sent();
                        return [2 /*return*/, new TokenResponse({
                                accessToken: response.data.access_token,
                                tokenType: response.data.token_type,
                                expiresIn: response.data.expires_in,
                                refreshToken: response.data.refresh_token,
                                scope: response.data.scope,
                                idToken: response.data.id_token,
                                issuedAt: response.data.issued_at
                            })];
                }
            });
        });
    };
    TokenRequest.prototype.getQueryBody = function () {
        var queryBody = new URLSearchParams();
        queryBody.set("grant_type", this.grantType);
        if (this.scopes) {
            queryBody.set("scope", this.scopes.join(" "));
        }
        queryBody.set("client_id", this.clientId);
        if (this.extraParams) {
            for (var extra in this.extraParams) {
                if (extra in this.extraParams && !(extra in queryBody)) {
                    queryBody.set(extra, this.extraParams[extra]);
                }
            }
        }
        return queryBody;
    };
    return TokenRequest;
}(Request));
/**
 * Access token request. Exchange an authorization code for a user access token.
 *
 * [Section 4.1.3](https://tools.ietf.org/html/rfc6749#section-4.1.3)
 */
var AccessTokenRequest = /** @class */ (function (_super) {
    __extends(AccessTokenRequest, _super);
    function AccessTokenRequest(options) {
        var _this = _super.call(this, options, types_1.GrantType.AuthorizationCode) || this;
        _this.code = options.code;
        _this.redirectUri = options.redirectUri;
        return _this;
    }
    AccessTokenRequest.prototype.getQueryBody = function () {
        var queryBody = _super.prototype.getQueryBody.call(this);
        if (this.redirectUri) {
            queryBody.set("redirect_uri", this.redirectUri);
        }
        if (this.code) {
            queryBody.set("code", this.code);
        }
        return queryBody;
    };
    AccessTokenRequest.prototype.getRequestConfig = function () {
        return {
            clientId: this.clientId,
            grantType: this.grantType,
            code: this.code,
            redirectUri: this.redirectUri,
            extraParams: this.extraParams,
            scopes: this.scopes
        };
    };
    return AccessTokenRequest;
}(TokenRequest));
exports.AccessTokenRequest = AccessTokenRequest;
/**
 * Refresh request.
 *
 * [Section 6](https://tools.ietf.org/html/rfc6749#section-6)
 */
var RefreshTokenRequest = /** @class */ (function (_super) {
    __extends(RefreshTokenRequest, _super);
    function RefreshTokenRequest(options) {
        var _this = _super.call(this, options, types_1.GrantType.RefreshToken) || this;
        _this.refreshToken = options.refreshToken;
        return _this;
    }
    RefreshTokenRequest.prototype.getQueryBody = function () {
        var queryBody = _super.prototype.getQueryBody.call(this);
        if (this.refreshToken) {
            queryBody.set("refresh_token", this.refreshToken);
        }
        return queryBody;
    };
    RefreshTokenRequest.prototype.getRequestConfig = function () {
        return {
            clientId: this.clientId,
            grantType: this.grantType,
            refreshToken: this.refreshToken,
            extraParams: this.extraParams,
            scopes: this.scopes
        };
    };
    return RefreshTokenRequest;
}(TokenRequest));
exports.RefreshTokenRequest = RefreshTokenRequest;
function exchangeCodeAsync(config, discovery) {
    var request = new AccessTokenRequest(config);
    console.log("ACCESS TOKEN EXCHANGE CONFIG", config);
    return request.performAsync(discovery);
}
exports.exchangeCodeAsync = exchangeCodeAsync;

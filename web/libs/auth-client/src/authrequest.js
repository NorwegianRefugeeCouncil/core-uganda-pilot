"use strict";
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
exports.AuthRequest = void 0;
var oauth_pkce_1 = require("oauth-pkce");
var browser_1 = require("./browser");
var types_1 = require("./types/types");
var tokenrequest_1 = require("./tokenrequest");
var queryparams_1 = require("./queryparams");
var error_1 = require("./error");
var _authLock = false;
var AuthRequest = /** @class */ (function () {
    function AuthRequest(request) {
        var _a, _b, _c, _d, _e;
        this.url = "";
        this.codeVerifier = "";
        this.codeChallenge = "";
        this.clientId = "";
        this.redirectUri = "";
        this.responseType = (_a = request.responseType) !== null && _a !== void 0 ? _a : types_1.ResponseType.Code;
        this.clientId = request.clientId;
        this.redirectUri = request.redirectUri;
        this.scopes = request.scopes ? request.scopes : [];
        this.prompt = request.prompt ? request.prompt : types_1.Prompt.Default;
        this.state = (_b = request.state) !== null && _b !== void 0 ? _b : (0, browser_1.generateRandom)(10);
        this.extraParams = (_c = request.extraParams) !== null && _c !== void 0 ? _c : {};
        this.codeChallengeMethod = (_d = request.codeChallengeMethod) !== null && _d !== void 0 ? _d : types_1.CodeChallengeMethod.S256;
        this.usePKCE = (_e = request.usePKCE) !== null && _e !== void 0 ? _e : true;
    }
    AuthRequest.prototype.promptAsync = function (discoveryDocument) {
        var _a;
        return __awaiter(this, void 0, void 0, function () {
            var url, result, ret;
            return __generator(this, function (_b) {
                switch (_b.label) {
                    case 0: return [4 /*yield*/, this.makeAuthUrlAsync(discoveryDocument)];
                    case 1:
                        url = _b.sent();
                        _b.label = 2;
                    case 2:
                        _b.trys.push([2, , 4, 5]);
                        _authLock = true;
                        return [4 /*yield*/, (0, browser_1.openAuthSessionAsync)({ url: url })];
                    case 3:
                        result = _b.sent();
                        return [3 /*break*/, 5];
                    case 4:
                        _authLock = false;
                        return [7 /*endfinally*/];
                    case 5:
                        if ((result === null || result === void 0 ? void 0 : result.type) !== types_1.WebBrowserResultType.SUCCESS) {
                            ret = { type: result === null || result === void 0 ? void 0 : result.type };
                            return [2 /*return*/, ret];
                        }
                        if ("url" in result) {
                            return [2 /*return*/, this.parseReturnUrl((_a = result) === null || _a === void 0 ? void 0 : _a.url)];
                        }
                        else {
                            throw new error_1.CodedError("ERR_NO_RESULT_URL", "url not found in auth result");
                        }
                        return [2 /*return*/];
                }
            });
        });
    };
    AuthRequest.prototype.parseReturnUrl = function (url) {
        var _a = (0, queryparams_1.getQueryParams)(url), params = _a.params, errorCode = _a.errorCode;
        var state = params.state, _b = params.error, error = _b === void 0 ? errorCode : _b;
        var parsedError = null;
        var authentication = null;
        if (state !== this.state) {
            // This is a non-standard error
            parsedError = new types_1.AuthError({
                error: 'state_mismatch',
                error_description: 'Cross-Site request verification failed. Cached state and returned state do not match.'
            });
        }
        else if (error) {
            parsedError = new types_1.AuthError(__assign({ error: error }, params));
        }
        if (params.access_token) {
            authentication = tokenrequest_1.TokenResponse.fromQueryParams(params);
        }
        var result = {
            type: parsedError ? 'error' : 'success',
            error: parsedError,
            url: url,
            params: params,
            authentication: authentication,
            // Return errorCode for legacy
            errorCode: errorCode
        };
        return result;
    };
    AuthRequest.prototype.makeAuthUrlAsync = function (discovery) {
        var _a;
        return __awaiter(this, void 0, void 0, function () {
            var request, params, extra, query;
            return __generator(this, function (_b) {
                switch (_b.label) {
                    case 0: return [4 /*yield*/, this.getAuthRequestConfigAsync()];
                    case 1:
                        request = _b.sent();
                        if (!request.state)
                            throw new Error('Cannot make request URL without a valid `state` loaded');
                        params = {};
                        if (request.codeChallenge) {
                            params.code_challenge = request.codeChallenge;
                        }
                        // copy over extra params
                        for (extra in request.extraParams) {
                            if (extra in request.extraParams) {
                                params[extra] = request.extraParams[extra];
                            }
                        }
                        if (request.usePKCE && request.codeChallengeMethod) {
                            params.code_challenge_method = request.codeChallengeMethod;
                        }
                        if (request.prompt) {
                            params.prompt = request.prompt;
                        }
                        // These overwrite any extra params
                        params.redirect_uri = request.redirectUri;
                        params.client_id = request.clientId;
                        params.response_type = request.responseType;
                        params.state = request.state;
                        if ((_a = request.scopes) === null || _a === void 0 ? void 0 : _a.length) {
                            params.scope = request.scopes.join(' ');
                        }
                        query = (0, queryparams_1.buildQueryString)(params);
                        // Store the URL for later
                        this.url = discovery.authorization_endpoint + "?" + query;
                        return [2 /*return*/, this.url];
                }
            });
        });
    };
    AuthRequest.prototype.ensureCodeIsSetupAsync = function () {
        return __awaiter(this, void 0, void 0, function () {
            var _this = this;
            return __generator(this, function (_a) {
                if (this.codeVerifier) {
                    return [2 /*return*/];
                }
                return [2 /*return*/, new Promise(function (resolve, reject) {
                        (0, oauth_pkce_1["default"])(43, function (error, value) {
                            if (error) {
                                reject(error);
                            }
                            _this.codeVerifier = value.verifier;
                            _this.codeChallenge = value.challenge;
                            resolve();
                        });
                    })];
            });
        });
    };
    AuthRequest.prototype.getAuthRequestConfigAsync = function () {
        return __awaiter(this, void 0, void 0, function () {
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        if (!this.usePKCE) return [3 /*break*/, 2];
                        return [4 /*yield*/, this.ensureCodeIsSetupAsync()];
                    case 1:
                        _a.sent();
                        _a.label = 2;
                    case 2: return [2 /*return*/, {
                            clientId: this.clientId,
                            responseType: this.responseType,
                            redirectUri: this.redirectUri,
                            scopes: this.scopes,
                            usePKCE: this.usePKCE,
                            codeChallengeMethod: this.codeChallengeMethod,
                            codeChallenge: this.codeChallenge,
                            prompt: this.prompt,
                            extraParams: this.extraParams,
                            state: this.state
                        }];
                }
            });
        });
    };
    return AuthRequest;
}());
exports.AuthRequest = AuthRequest;

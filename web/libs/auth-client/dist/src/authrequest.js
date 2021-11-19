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
exports.AuthRequest = void 0;
const oauth_pkce_1 = __importDefault(require("oauth-pkce"));
const browser_1 = require("./browser");
const types_1 = require("./types/types");
const tokenrequest_1 = require("./tokenrequest");
const queryparams_1 = require("./queryparams");
const error_1 = require("./error");
let _authLock = false;
class AuthRequest {
    constructor(request) {
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
    promptAsync(discoveryDocument) {
        var _a;
        return __awaiter(this, void 0, void 0, function* () {
            const url = yield this.makeAuthUrlAsync(discoveryDocument);
            let result;
            try {
                _authLock = true;
                result = yield (0, browser_1.openAuthSessionAsync)({ url });
            }
            finally {
                _authLock = false;
            }
            if ((result === null || result === void 0 ? void 0 : result.type) !== types_1.WebBrowserResultType.SUCCESS) {
                const ret = { type: result === null || result === void 0 ? void 0 : result.type };
                return ret;
            }
            if ("url" in result) {
                return this.parseReturnUrl((_a = result) === null || _a === void 0 ? void 0 : _a.url);
            }
            else {
                throw new error_1.CodedError("ERR_NO_RESULT_URL", "url not found in auth result");
            }
        });
    }
    parseReturnUrl(url) {
        const { params, errorCode } = (0, queryparams_1.getQueryParams)(url);
        const { state, error = errorCode } = params;
        let parsedError = null;
        let authentication = null;
        if (state !== this.state) {
            // This is a non-standard error
            parsedError = new types_1.AuthError({
                error: 'state_mismatch',
                error_description: 'Cross-Site request verification failed. Cached state and returned state do not match.',
            });
        }
        else if (error) {
            parsedError = new types_1.AuthError(Object.assign({ error }, params));
        }
        if (params.access_token) {
            authentication = tokenrequest_1.TokenResponse.fromQueryParams(params);
        }
        const result = {
            type: parsedError ? 'error' : 'success',
            error: parsedError,
            url,
            params,
            authentication,
            // Return errorCode for legacy
            errorCode,
        };
        return result;
    }
    makeAuthUrlAsync(discovery) {
        var _a;
        return __awaiter(this, void 0, void 0, function* () {
            const request = yield this.getAuthRequestConfigAsync();
            if (!request.state)
                throw new Error('Cannot make request URL without a valid `state` loaded');
            // Create a query string
            const params = {};
            if (request.codeChallenge) {
                params.code_challenge = request.codeChallenge;
            }
            // copy over extra params
            for (const extra in request.extraParams) {
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
            const query = (0, queryparams_1.buildQueryString)(params);
            // Store the URL for later
            this.url = `${discovery.authorization_endpoint}?${query}`;
            return this.url;
        });
    }
    ensureCodeIsSetupAsync() {
        return __awaiter(this, void 0, void 0, function* () {
            if (this.codeVerifier) {
                return;
            }
            return new Promise((resolve, reject) => {
                (0, oauth_pkce_1.default)(43, (error, value) => {
                    if (error) {
                        reject(error);
                    }
                    this.codeVerifier = value.verifier;
                    this.codeChallenge = value.challenge;
                    resolve();
                });
            });
        });
    }
    getAuthRequestConfigAsync() {
        return __awaiter(this, void 0, void 0, function* () {
            if (this.usePKCE) {
                yield this.ensureCodeIsSetupAsync();
            }
            return {
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
            };
        });
    }
}
exports.AuthRequest = AuthRequest;

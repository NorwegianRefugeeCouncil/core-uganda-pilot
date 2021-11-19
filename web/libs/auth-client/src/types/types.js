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
exports.__esModule = true;
exports.WebBrowserResultType = exports.GrantType = exports.TokenTypeHint = exports.Prompt = exports.CodeChallengeMethod = exports.ResponseType = exports.TokenError = exports.AuthError = exports.ResponseError = void 0;
var error_1 = require("../error");
var errorCodeMessages = {
    // https://tools.ietf.org/html/rfc6749#section-4.1.2.1
    // https://openid.net/specs/openid-connect-core-1_0.html#AuthError
    auth: {
        // OAuth 2.0
        invalid_request: "The request is missing a required parameter, includes an invalid parameter value, includes a parameter more than once, or is otherwise malformed.",
        unauthorized_client: "The client is not authorized to request an authorization code using this method.",
        access_denied: "The resource owner or authorization server denied the request.",
        unsupported_response_type: "The authorization server does not support obtaining an authorization code using this method.",
        invalid_scope: 'The requested scope is invalid, unknown, or malformed.',
        server_error: 'The authorization server encountered an unexpected condition that prevented it from fulfilling the request. (This error code is needed because a 500 Internal Server Error HTTP status code cannot be returned to the client via an HTTP redirect.)',
        temporarily_unavailable: 'The authorization server is currently unable to handle the request due to a temporary overloading or maintenance of the server.  (This error code is needed because a 503 Service Unavailable HTTP status code cannot be returned to the client via an HTTP redirect.)',
        // Open ID Connect error codes
        interaction_required: 'Auth server requires user interaction of some form to proceed. This error may be returned when the prompt parameter value in the auth request is none, but the auth request cannot be completed without displaying a user interface for user interaction.',
        login_required: 'Auth server requires user authentication. This error may be returned when the prompt parameter value in the auth request is none, but the auth request cannot be completed without displaying a user interface for user authentication.',
        account_selection_required: 'User is required to select a session at the auth server. The user may be authenticated at the auth server with different associated accounts, but the user did not select a session. This error may be returned when the prompt parameter value in the auth request is `none`, but the auth request cannot be completed without displaying a user interface to prompt for a session to use.',
        consent_required: 'Auth server requires user consent. This error may be returned when the prompt parameter value in the auth request is none, but the auth request cannot be completed without displaying a user interface for user consent.',
        invalid_request_uri: 'The `request_uri` in the auth request returns an error or contains invalid data.',
        invalid_request_object: 'The request parameter contains an invalid request object.',
        request_not_supported: 'The OP does not support use of the `request` parameter defined in Section 6. (https://openid.net/specs/openid-connect-core-1_0.html#JWTRequests)',
        request_uri_not_supported: 'The OP does not support use of the `request_uri` parameter defined in Section 6. (https://openid.net/specs/openid-connect-core-1_0.html#JWTRequests)',
        registration_not_supported: 'The OP does not support use of the `registration` parameter defined in Section 7.2.1. (https://openid.net/specs/openid-connect-core-1_0.html#RegistrationParameter)'
    },
    // https://tools.ietf.org/html/rfc6749#section-5.2
    token: {
        invalid_request: "The request is missing a required parameter, includes an unsupported parameter value (other than grant type), repeats a parameter, includes multiple credentials, utilizes more than one mechanism for authenticating the client, or is otherwise malformed.",
        invalid_client: "Client authentication failed (e.g., unknown client, no client authentication included, or unsupported authentication method).  The authorization server MAY return an HTTP 401 (Unauthorized) status code to indicate which HTTP authentication schemes are supported.  If the client attempted to authenticate via the \"Authorization\" request header field, the authorization server MUST respond with an HTTP 401 (Unauthorized) status code and include the \"WWW-Authenticate\" response header field matching the authentication scheme used by the client.",
        invalid_grant: "The provided authorization grant (e.g., authorization code, resource owner credentials) or refresh token is invalid, expired, revoked, does not match the redirection URI used in the authorization request, or was issued to another client.",
        unauthorized_client: "The authenticated client is not authorized to use this authorization grant type.",
        unsupported_grant_type: "The authorization grant type is not supported by the authorization server."
    }
};
var ResponseError = /** @class */ (function (_super) {
    __extends(ResponseError, _super);
    function ResponseError(params, errorCodeType) {
        var _this = this;
        var error = params.error, error_description = params.error_description, error_uri = params.error_uri;
        var message = errorCodeMessages[errorCodeType][error];
        var errorMessage;
        if (message) {
            errorMessage = message + (error_description ? "\nMore info: " + error_description : '');
        }
        else if (error_description) {
            errorMessage = error_description;
        }
        else {
            errorMessage = 'An unknown error occurred';
        }
        _this = _super.call(this, error, errorMessage) || this;
        _this.description = error_description !== null && error_description !== void 0 ? error_description : message;
        _this.uri = error_uri;
        _this.params = params;
        return _this;
    }
    return ResponseError;
}(error_1.CodedError));
exports.ResponseError = ResponseError;
/**
 * [Section 5.2](https://tools.ietf.org/html/rfc6749#section-5.2)
 */
var AuthError = /** @class */ (function (_super) {
    __extends(AuthError, _super);
    function AuthError(response) {
        var _this = _super.call(this, response, 'auth') || this;
        _this.state = response.state;
        return _this;
    }
    return AuthError;
}(ResponseError));
exports.AuthError = AuthError;
/**
 * [Section 4.1.2.1](https://tools.ietf.org/html/rfc6749#section-4.1.2.1)
 */
var TokenError = /** @class */ (function (_super) {
    __extends(TokenError, _super);
    function TokenError(response) {
        return _super.call(this, response, 'token') || this;
    }
    return TokenError;
}(ResponseError));
exports.TokenError = TokenError;
var ResponseType;
(function (ResponseType) {
    ResponseType["Code"] = "code";
    ResponseType["Token"] = "token";
    ResponseType["IdToken"] = "id_token";
})(ResponseType = exports.ResponseType || (exports.ResponseType = {}));
var CodeChallengeMethod;
(function (CodeChallengeMethod) {
    CodeChallengeMethod["S256"] = "S256";
    CodeChallengeMethod["Plain"] = "plain";
})(CodeChallengeMethod = exports.CodeChallengeMethod || (exports.CodeChallengeMethod = {}));
var Prompt;
(function (Prompt) {
    Prompt["None"] = "none";
    Prompt["Login"] = "login";
    Prompt["Consent"] = "consent";
    Prompt["SelectAccount"] = "select_account";
    Prompt["Default"] = "";
})(Prompt = exports.Prompt || (exports.Prompt = {}));
/**
 * A hint about the type of the token submitted for revocation.
 *
 * [Section 2.1](https://tools.ietf.org/html/rfc7009#section-2.1)
 */
var TokenTypeHint;
(function (TokenTypeHint) {
    /**
     * Access token.
     *
     * [Section 1.4](https://tools.ietf.org/html/rfc6749#section-1.4)
     */
    TokenTypeHint["AccessToken"] = "access_token";
    /**
     * Refresh token.
     *
     * [Section 1.5](https://tools.ietf.org/html/rfc6749#section-1.5)
     */
    TokenTypeHint["RefreshToken"] = "refresh_token";
})(TokenTypeHint = exports.TokenTypeHint || (exports.TokenTypeHint = {}));
/**
 * Grant type values used in dynamic client registration and auth requests.
 *
 * [Appendix A.10](https://tools.ietf.org/html/rfc6749#appendix-A.10)
 */
var GrantType;
(function (GrantType) {
    /**
     * Used for exchanging an authorization code for one or more tokens.
     *
     * [Section 4.1.3](https://tools.ietf.org/html/rfc6749#section-4.1.3)
     */
    GrantType["AuthorizationCode"] = "authorization_code";
    /**
     * Used when obtaining an access token.
     *
     * [Section 4.2](https://tools.ietf.org/html/rfc6749#section-4.2)
     */
    GrantType["Implicit"] = "implicit";
    /**
     * Used when exchanging a refresh token for a new token.
     *
     * [Section 6](https://tools.ietf.org/html/rfc6749#section-6)
     */
    GrantType["RefreshToken"] = "refresh_token";
    /**
     * Used for client credentials flow.
     *
     * [Section 4.4.2](https://tools.ietf.org/html/rfc6749#section-4.4.2)
     */
    GrantType["ClientCredentials"] = "client_credentials";
})(GrantType = exports.GrantType || (exports.GrantType = {}));
var WebBrowserResultType;
(function (WebBrowserResultType) {
    WebBrowserResultType["CANCEL"] = "cancel";
    WebBrowserResultType["DISMISS"] = "dismiss";
    WebBrowserResultType["OPENED"] = "opened";
    WebBrowserResultType["LOCKED"] = "locked";
    WebBrowserResultType["SUCCESS"] = "success";
})(WebBrowserResultType = exports.WebBrowserResultType || (exports.WebBrowserResultType = {}));

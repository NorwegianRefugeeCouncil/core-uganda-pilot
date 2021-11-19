"use strict";
exports.__esModule = true;
exports.IdentityProvider = exports.Organization = exports.OAuth2Client = void 0;
var OAuth2Client = /** @class */ (function () {
    function OAuth2Client() {
        this.id = "";
        this.clientName = "";
        this.clientSecret = "";
        this.uri = "";
        this.grantTypes = ["authorization_code"];
        this.responseTypes = ["code"];
        this.scope = "";
        this.redirectUris = [];
        this.allowedCorsOrigins = [];
        this.tokenEndpointAuthMethod = "client_secret_basic";
    }
    return OAuth2Client;
}());
exports.OAuth2Client = OAuth2Client;
var Organization = /** @class */ (function () {
    function Organization() {
        this.id = "";
        this.key = "";
        this.name = "";
    }
    return Organization;
}());
exports.Organization = Organization;
var IdentityProvider = /** @class */ (function () {
    function IdentityProvider() {
        this.id = "";
        this.name = "";
        this.organizationId = "";
        this.domain = "";
        this.clientId = "";
        this.clientSecret = "";
        this.emailDomain = "";
    }
    return IdentityProvider;
}());
exports.IdentityProvider = IdentityProvider;

"use strict";
exports.__esModule = true;
var rxjs_1 = require("rxjs");
var operators_1 = require("rxjs/operators");
var models_1 = require("./types/models");
var ajax_1 = require("rxjs/ajax");
var xhr2_1 = require("xhr2");
// needed for rxjs/ajax compatibility outside the browser
global.XMLHttpRequest = global.XMLHttpRequest ? global.XMLHttpRequest : xhr2_1.XMLHttpRequest;
// todo: should come from environment
var shouldAddAuthHeader = true;
var noop = function () { };
var HttpClient = /** @class */ (function () {
    function HttpClient(shouldAddAuthPassthroughHeader) {
        this.headers = {};
        if (shouldAddAuthPassthroughHeader) {
            this.headers = {
                'X-Authenticated-User-Subject': 'stephen.kabagambe@email.com'
            };
        }
    }
    HttpClient.prototype.get = function (url) {
        return (0, ajax_1.ajax)({
            url: url,
            headers: this.headers,
            method: 'GET',
            async: true,
            timeout: 0,
            crossDomain: true,
            withCredentials: false,
            responseType: 'json'
        });
    };
    HttpClient.prototype.put = function (url, body) {
        return (0, ajax_1.ajax)({
            url: url,
            body: body,
            headers: this.headers,
            method: 'PUT',
            async: true,
            timeout: 0,
            crossDomain: true,
            withCredentials: false,
            responseType: 'json'
        });
    };
    HttpClient.prototype.post = function (url, body) {
        return (0, ajax_1.ajax)({
            url: url,
            body: body,
            headers: this.headers,
            method: 'POST',
            async: true,
            timeout: 0,
            crossDomain: true,
            withCredentials: false,
            responseType: 'json'
        });
    };
    HttpClient.prototype["delete"] = function (url) {
        return (0, ajax_1.ajax)({
            url: url,
            headers: this.headers,
            method: 'DELETE',
            async: true,
            timeout: 0,
            crossDomain: true,
            withCredentials: false,
            responseType: 'json'
        });
    };
    return HttpClient;
}());
var CaseClient = /** @class */ (function () {
    function CaseClient() {
        this.httpClient = new HttpClient(shouldAddAuthHeader);
        this.endpoint = 'http://localhost:9000/apis/cms/v1/cases';
    }
    CaseClient.prototype.Get = function (id) {
        return this.httpClient.get(this.endpoint + id);
    };
    CaseClient.prototype.Create = function (c) {
        return this.httpClient.post(this.endpoint, c);
    };
    CaseClient.prototype.Update = function (c) {
        return this.httpClient.put(this.endpoint + c.id, c);
    };
    CaseClient.prototype.List = function (lo) {
        var query = new URLSearchParams(lo);
        return this.httpClient.get(query ? this.endpoint : this.endpoint + ("?" + query));
    };
    return CaseClient;
}());
var CaseTypeClient = /** @class */ (function () {
    function CaseTypeClient() {
        this.httpClient = new HttpClient(shouldAddAuthHeader);
        this.endpoint = 'http://localhost:9000/apis/cms/v1/casetypes';
    }
    CaseTypeClient.prototype.Get = function (id) {
        return this.httpClient.get(this.endpoint + id);
    };
    CaseTypeClient.prototype.Create = function (c) {
        return this.httpClient.post(this.endpoint, c);
    };
    CaseTypeClient.prototype.Update = function (c) {
        return this.httpClient.put(this.endpoint + c.id, c);
    };
    CaseTypeClient.prototype.List = function (lo) {
        var query = new URLSearchParams(lo);
        return this.httpClient.get(query ? this.endpoint : this.endpoint + ("?" + query));
    };
    return CaseTypeClient;
}());
var CommentClient = /** @class */ (function () {
    function CommentClient() {
        this.httpClient = new HttpClient(shouldAddAuthHeader);
        this.endpoint = 'http://localhost:9000/apis/cms/v1/comments';
    }
    CommentClient.prototype.Get = function (id) {
        return this.httpClient.get(this.endpoint + id);
    };
    CommentClient.prototype.Create = function (c) {
        return this.httpClient.post(this.endpoint, c);
    };
    CommentClient.prototype.Update = function (c) {
        return this.httpClient.put(this.endpoint + c.id, c);
    };
    CommentClient.prototype.List = function (lo) {
        var query = new URLSearchParams(lo);
        return this.httpClient.get(query ? this.endpoint : this.endpoint + ("?" + query));
    };
    CommentClient.prototype.Delete = function (id) {
        return this.httpClient["delete"](this.endpoint + id);
    };
    return CommentClient;
}());
var CMSClient = /** @class */ (function () {
    function CMSClient() {
    }
    CMSClient.Cases = function () {
        return new CaseClient;
    };
    CMSClient.CaseTypes = function () {
        return new CaseTypeClient;
    };
    CMSClient.Comments = function () {
        return new CommentClient;
    };
    return CMSClient;
}());
CMSClient.Comments().List(new models_1.CommentListOptions())
    .pipe((0, operators_1.map)(function (response) {
    console.log(response);
}, (0, operators_1.catchError)(function (error) {
    console.log('error: ', error);
    return (0, rxjs_1.of)(error);
}))).subscribe(noop);

"use strict";
exports.__esModule = true;
var react_1 = require("react");
var browser_1 = require("../browser");
var tokenrequest_1 = require("../tokenrequest");
var types_1 = require("../types/types");
var hooks_1 = require("../hooks");
var axios_1 = require("axios");
(0, browser_1.maybeCompleteAuthSession)();
// TODO: https://betterprogramming.pub/building-secure-login-flow-with-oauth-2-openid-in-react-apps-ce6e8e29630a
var AuthWrapper = function (_a) {
    var children = _a.children, _b = _a.scopes, scopes = _b === void 0 ? [] : _b, clientId = _a.clientId, _c = _a.axiosInstance, axiosInstance = _c === void 0 ? axios_1["default"].create() : _c, issuer = _a.issuer, _d = _a.redirectUriSuffix, redirectUriSuffix = _d === void 0 ? '/' : _d, customLoginComponent = _a.customLoginComponent, _e = _a.handleLoginErr, handleLoginErr = _e === void 0 ? console.log : _e;
    var redirectUri = window.location.protocol + "//" + window.location.host + "/" + redirectUriSuffix;
    var discovery = (0, hooks_1.useDiscovery)(issuer);
    var _f = react_1["default"].useState(), tokenResponse = _f[0], setTokenResponse = _f[1];
    var _g = react_1["default"].useState(false), isLoggedIn = _g[0], setIsLoggedIn = _g[1];
    var _h = (0, hooks_1.useAuthRequest)({
        clientId: clientId,
        usePKCE: true,
        responseType: types_1.ResponseType.Code,
        codeChallengeMethod: types_1.CodeChallengeMethod.S256,
        scopes: scopes,
        redirectUri: redirectUri
    }, discovery), request = _h[0], response = _h[1], promptAsync = _h[2];
    react_1["default"].useEffect(function () {
        if (!discovery) {
            return;
        }
        if (!(request === null || request === void 0 ? void 0 : request.codeVerifier)) {
            return;
        }
        if (!response || (response === null || response === void 0 ? void 0 : response.type) !== "success") {
            return;
        }
        var exchangeConfig = {
            code: response.params.code,
            clientId: clientId,
            redirectUri: redirectUri,
            extraParams: {
                "code_verifier": request === null || request === void 0 ? void 0 : request.codeVerifier
            }
        };
        (0, tokenrequest_1.exchangeCodeAsync)(exchangeConfig, discovery)
            .then(function (a) {
            setTokenResponse(a);
        })["catch"](function (err) {
            setTokenResponse(undefined);
        });
    }, [request === null || request === void 0 ? void 0 : request.codeVerifier, response, discovery]);
    react_1["default"].useEffect(function () {
        if (!discovery) {
            return;
        }
        if (tokenResponse === null || tokenResponse === void 0 ? void 0 : tokenResponse.shouldRefresh()) {
            var refreshConfig = {
                clientId: clientId,
                scopes: scopes,
                extraParams: {}
            };
            tokenResponse === null || tokenResponse === void 0 ? void 0 : tokenResponse.refreshAsync(refreshConfig, discovery).then(function (resp) {
                setTokenResponse(resp);
            })["catch"](function (err) {
                setTokenResponse(undefined);
            });
        }
    }, [tokenResponse === null || tokenResponse === void 0 ? void 0 : tokenResponse.shouldRefresh(), discovery]);
    react_1["default"].useEffect(function () {
        if (tokenResponse) {
            if (!isLoggedIn) {
                setIsLoggedIn(true);
            }
        }
        else {
            if (isLoggedIn) {
                setIsLoggedIn(false);
            }
        }
    }, [tokenResponse, isLoggedIn]);
    react_1["default"].useEffect(function () {
        var interceptor = axiosInstance.interceptors.request.use(function (value) {
            if (!(tokenResponse === null || tokenResponse === void 0 ? void 0 : tokenResponse.accessToken)) {
                return value;
            }
            if (!value.headers) {
                value.headers = {};
            }
            value.headers["Authorization"] = "Bearer " + tokenResponse.accessToken;
            return value;
        });
        return function () {
            axiosInstance.interceptors.request.eject(interceptor);
        };
    }, [tokenResponse === null || tokenResponse === void 0 ? void 0 : tokenResponse.accessToken]);
    var handleLogin = react_1["default"].useCallback(function () {
        promptAsync()["catch"](function (err) {
            handleLoginErr(err);
        });
    }, [discovery, request, promptAsync]);
    if (!isLoggedIn) {
        return (<react_1["default"].Fragment>
                {customLoginComponent
                ?
                    customLoginComponent({ login: handleLogin })
                :
                    <button onClick={handleLogin}>Login</button>}
            </react_1["default"].Fragment>);
    }
    return (<react_1["default"].Fragment>
            {children}
        </react_1["default"].Fragment>);
};
exports["default"] = AuthWrapper;

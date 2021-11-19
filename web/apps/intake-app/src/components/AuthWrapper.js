"use strict";
exports.__esModule = true;
exports.AuthWrapper = void 0;
var react_1 = require("react");
var clients_1 = require("../utils/clients");
var expo_auth_session_1 = require("expo-auth-session");
var react_native_1 = require("react-native");
var expo_constants_1 = require("expo-constants");
var AuthWrapper = function (props) {
    var _a, _b, _c, _d, _e, _f;
    var children = props.children;
    var clientId = (_b = (_a = expo_constants_1["default"].manifest) === null || _a === void 0 ? void 0 : _a.extra) === null || _b === void 0 ? void 0 : _b.client_id;
    var useProxy = (0, react_1.useMemo)(function () { return react_native_1.Platform.select({ web: false, "default": false }); }, []);
    var redirectUri = (0, react_1.useMemo)(function () { return (0, expo_auth_session_1.makeRedirectUri)({ scheme: 'nrccore' }); }, []);
    var discovery = (0, expo_auth_session_1.useAutoDiscovery)((_d = (_c = expo_constants_1["default"].manifest) === null || _c === void 0 ? void 0 : _c.extra) === null || _d === void 0 ? void 0 : _d.issuer);
    var _g = (0, react_1.useState)(false), loggedIn = _g[0], setLoggedIn = _g[1];
    var _h = (0, react_1.useState)(), tokenResponse = _h[0], setTokenResponse = _h[1];
    var _j = (0, expo_auth_session_1.useAuthRequest)({
        clientId: clientId,
        usePKCE: true,
        responseType: expo_auth_session_1.ResponseType.Code,
        codeChallengeMethod: expo_auth_session_1.CodeChallengeMethod.S256,
        scopes: (_f = (_e = expo_constants_1["default"].manifest) === null || _e === void 0 ? void 0 : _e.extra) === null || _f === void 0 ? void 0 : _f.scopes,
        redirectUri: redirectUri
    }, discovery), request = _j[0], response = _j[1], promptAsync = _j[2];
    react_1["default"].useEffect(function () {
        if (!discovery) {
            return;
        }
        if (!(request === null || request === void 0 ? void 0 : request.codeVerifier)) {
            return;
        }
        if (!response || response.type !== "success") {
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
        (0, expo_auth_session_1.exchangeCodeAsync)(exchangeConfig, discovery)
            .then(function (a) {
            setTokenResponse(a);
        })["catch"](function (err) {
            setTokenResponse(undefined);
        });
    }, [request === null || request === void 0 ? void 0 : request.codeVerifier, response, discovery]);
    (0, react_1.useEffect)(function () {
        var _a, _b;
        if (!discovery) {
            return;
        }
        if (tokenResponse === null || tokenResponse === void 0 ? void 0 : tokenResponse.shouldRefresh()) {
            var refreshConfig = {
                clientId: clientId,
                scopes: (_b = (_a = expo_constants_1["default"].manifest) === null || _a === void 0 ? void 0 : _a.extra) === null || _b === void 0 ? void 0 : _b.scopes,
                extraParams: {}
            };
            tokenResponse === null || tokenResponse === void 0 ? void 0 : tokenResponse.refreshAsync(refreshConfig, discovery).then(function (resp) {
                setTokenResponse(resp);
            })["catch"](function (err) {
                setTokenResponse(undefined);
            });
        }
    }, [tokenResponse === null || tokenResponse === void 0 ? void 0 : tokenResponse.shouldRefresh(), discovery]);
    (0, react_1.useEffect)(function () {
        if (tokenResponse) {
            if (!loggedIn) {
                setLoggedIn(true);
            }
        }
        else {
            if (loggedIn) {
                setLoggedIn(false);
            }
        }
    }, [tokenResponse, loggedIn]);
    (0, react_1.useEffect)(function () {
        var interceptor = clients_1.axiosInstance.interceptors.request.use(function (value) {
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
            clients_1.axiosInstance.interceptors.request.eject(interceptor);
        };
    }, [tokenResponse === null || tokenResponse === void 0 ? void 0 : tokenResponse.accessToken]);
    var handleLogin = (0, react_1.useCallback)(function () {
        promptAsync({ useProxy: useProxy }).then(function (response) {
            console.log("PROMPT RESPONSE", response);
        })["catch"](function (err) {
            console.log("PROMPT ERROR", err);
        });
    }, [useProxy, promptAsync]);
    if (!loggedIn) {
        return (<react_native_1.Button title={"Login"} disabled={!request} onPress={handleLogin}/>);
    }
    return (<react_1.Fragment>
            {children}
        </react_1.Fragment>);
};
exports.AuthWrapper = AuthWrapper;

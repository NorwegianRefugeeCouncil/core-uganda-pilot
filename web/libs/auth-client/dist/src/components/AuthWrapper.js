"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const jsx_runtime_1 = require("react/jsx-runtime");
const react_1 = __importDefault(require("react"));
const browser_1 = require("../browser");
const tokenrequest_1 = require("../tokenrequest");
const types_1 = require("../types/types");
const hooks_1 = require("../hooks");
const axios_1 = __importDefault(require("axios"));
(0, browser_1.maybeCompleteAuthSession)();
// TODO: https://betterprogramming.pub/building-secure-login-flow-with-oauth-2-openid-in-react-apps-ce6e8e29630a
const AuthWrapper = ({ children, scopes = [], clientId, axiosInstance = axios_1.default.create(), issuer, redirectUriSuffix = '/', customLoginComponent, handleLoginErr = console.log, }) => {
    const redirectUri = `${window.location.protocol}//${window.location.host}/${redirectUriSuffix}`;
    const discovery = (0, hooks_1.useDiscovery)(issuer);
    const [tokenResponse, setTokenResponse] = react_1.default.useState();
    const [isLoggedIn, setIsLoggedIn] = react_1.default.useState(false);
    const [request, response, promptAsync] = (0, hooks_1.useAuthRequest)({
        clientId,
        usePKCE: true,
        responseType: types_1.ResponseType.Code,
        codeChallengeMethod: types_1.CodeChallengeMethod.S256,
        scopes,
        redirectUri
    }, discovery);
    react_1.default.useEffect(() => {
        if (!discovery) {
            return;
        }
        if (!(request === null || request === void 0 ? void 0 : request.codeVerifier)) {
            return;
        }
        if (!response || (response === null || response === void 0 ? void 0 : response.type) !== "success") {
            return;
        }
        const exchangeConfig = {
            code: response.params.code,
            clientId,
            redirectUri,
            extraParams: {
                "code_verifier": request === null || request === void 0 ? void 0 : request.codeVerifier,
            }
        };
        (0, tokenrequest_1.exchangeCodeAsync)(exchangeConfig, discovery)
            .then(a => {
            setTokenResponse(a);
        })
            .catch((err) => {
            setTokenResponse(undefined);
        });
    }, [request === null || request === void 0 ? void 0 : request.codeVerifier, response, discovery]);
    react_1.default.useEffect(() => {
        if (!discovery) {
            return;
        }
        if (tokenResponse === null || tokenResponse === void 0 ? void 0 : tokenResponse.shouldRefresh()) {
            const refreshConfig = {
                clientId,
                scopes,
                extraParams: {}
            };
            tokenResponse === null || tokenResponse === void 0 ? void 0 : tokenResponse.refreshAsync(refreshConfig, discovery).then(resp => {
                setTokenResponse(resp);
            }).catch((err) => {
                setTokenResponse(undefined);
            });
        }
    }, [tokenResponse === null || tokenResponse === void 0 ? void 0 : tokenResponse.shouldRefresh(), discovery]);
    react_1.default.useEffect(() => {
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
    react_1.default.useEffect(() => {
        const interceptor = axiosInstance.interceptors.request.use(value => {
            if (!(tokenResponse === null || tokenResponse === void 0 ? void 0 : tokenResponse.accessToken)) {
                return value;
            }
            if (!value.headers) {
                value.headers = {};
            }
            value.headers["Authorization"] = `Bearer ${tokenResponse.accessToken}`;
            return value;
        });
        return () => {
            axiosInstance.interceptors.request.eject(interceptor);
        };
    }, [tokenResponse === null || tokenResponse === void 0 ? void 0 : tokenResponse.accessToken]);
    const handleLogin = react_1.default.useCallback(() => {
        promptAsync().catch((err) => {
            handleLoginErr(err);
        });
    }, [discovery, request, promptAsync]);
    if (!isLoggedIn) {
        return ((0, jsx_runtime_1.jsx)(react_1.default.Fragment, { children: customLoginComponent
                ?
                    customLoginComponent({ login: handleLogin })
                :
                    (0, jsx_runtime_1.jsx)("button", Object.assign({ onClick: handleLogin }, { children: "Login" }), void 0) }, void 0));
    }
    return ((0, jsx_runtime_1.jsx)(react_1.default.Fragment, { children: children }, void 0));
};
exports.default = AuthWrapper;

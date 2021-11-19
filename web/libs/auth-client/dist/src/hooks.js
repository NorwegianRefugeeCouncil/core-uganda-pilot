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
exports.useAuthRequest = exports.useAuthRequestResult = exports.useLoadedAuthRequest = exports.useDiscovery = void 0;
const react_1 = __importDefault(require("react"));
const discovery_1 = require("./discovery");
const authrequest_1 = require("./authrequest");
function useDiscovery(issuerOrDiscovery) {
    const [discoveryDocument, setDiscoveryDocument] = react_1.default.useState(null);
    react_1.default.useEffect(() => {
        let isAllowed = true;
        (0, discovery_1.resolveDiscoveryAsync)(issuerOrDiscovery).then(discovery => {
            if (isAllowed) {
                setDiscoveryDocument(discovery);
            }
        });
        return () => {
            isAllowed = false;
        };
    }, [issuerOrDiscovery]);
    return discoveryDocument;
}
exports.useDiscovery = useDiscovery;
function useLoadedAuthRequest(config, discovery, AuthRequestInstance) {
    const [request, setRequest] = react_1.default.useState(null);
    const scopeString = react_1.default.useMemo(() => { var _a; return (_a = config.scopes) === null || _a === void 0 ? void 0 : _a.join(','); }, [config.scopes]);
    const extraParamsString = react_1.default.useMemo(() => JSON.stringify(config.extraParams || {}), [config.extraParams]);
    react_1.default.useEffect(() => {
        let isMounted = true;
        if (discovery) {
            const request = new AuthRequestInstance(config);
            request.makeAuthUrlAsync(discovery).then(() => {
                if (isMounted) {
                    setRequest(request);
                }
            });
        }
        return () => {
            isMounted = false;
        };
    }, 
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [
        discovery === null || discovery === void 0 ? void 0 : discovery.authorization_endpoint,
        config.clientId,
        config.redirectUri,
        config.responseType,
        config.prompt,
        config.codeChallenge,
        config.state,
        config.usePKCE,
        scopeString,
        extraParamsString,
    ]);
    return request;
}
exports.useLoadedAuthRequest = useLoadedAuthRequest;
function useAuthRequestResult(request, discovery) {
    const [result, setResult] = react_1.default.useState(null);
    const promptAsync = react_1.default.useCallback(() => __awaiter(this, void 0, void 0, function* () {
        if (!discovery || !request) {
            throw new Error('Cannot prompt to authenticate until the request has finished loading.');
        }
        const result = yield (request === null || request === void 0 ? void 0 : request.promptAsync(discovery));
        setResult(result);
        return result;
    }), 
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [request === null || request === void 0 ? void 0 : request.url, discovery === null || discovery === void 0 ? void 0 : discovery.authorization_endpoint]);
    return [result, promptAsync];
}
exports.useAuthRequestResult = useAuthRequestResult;
function useAuthRequest(config, discovery) {
    const request = useLoadedAuthRequest(config, discovery, authrequest_1.AuthRequest);
    const [result, promptAsync] = useAuthRequestResult(request, discovery);
    return [request, result, promptAsync];
}
exports.useAuthRequest = useAuthRequest;

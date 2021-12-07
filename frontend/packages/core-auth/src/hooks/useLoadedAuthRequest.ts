import {AuthRequestConfig, DiscoveryDocument} from "../types/types";
import {AuthRequest} from "../types/authrequest";
import React from "react";
import Browser from "../types/browser";

export default function useLoadedAuthRequest(
    config: AuthRequestConfig,
    discovery: DiscoveryDocument | null,
    AuthRequestInstance: typeof AuthRequest,
    browser: Browser
): AuthRequest | null {
    const [authRequest, setAuthRequest] = React.useState<AuthRequest | null>(null);
    const scopeString = React.useMemo(() => config.scopes?.join(','), [config.scopes]);
    const extraParamsString = React.useMemo(
        () => JSON.stringify(config.extraParams || {}),
        [config.extraParams]
    );
    React.useEffect(
        () => {
            let isMounted = true;
            if (discovery) {
                const request = new AuthRequestInstance(config, browser);
                request.makeAuthUrlAsync(discovery).then(() => {
                    if (isMounted) {
                        setAuthRequest(request);
                    }
                });
            }
            return () => {
                isMounted = false;
            };
        },
        // eslint-disable-next-line react-hooks/exhaustive-deps
        [
            discovery?.authorization_endpoint,
            config.clientId,
            config.redirectUri,
            config.responseType,
            config.prompt,
            config.codeChallenge,
            config.state,
            config.usePKCE,
            scopeString,
            extraParamsString,
        ]
    );
    return authRequest;
}



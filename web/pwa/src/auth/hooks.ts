import {useCallback, useEffect, useMemo, useState} from "react";
import {DiscoveryDocument, IssuerOrDiscovery, resolveDiscoveryAsync} from "./discovery";
import {AuthRequestConfig, AuthRequestPromptOptions, AuthSessionResult} from "./types";
import {AuthRequest} from "./authrequest";

export function useDiscovery(issuerOrDiscovery: IssuerOrDiscovery): DiscoveryDocument | null {
    const [discoveryDocument, setDiscoveryDocument] = useState<DiscoveryDocument | null>(null)
    useEffect(() => {
        let isAllowed = true
        resolveDiscoveryAsync(issuerOrDiscovery).then(discovery => {
            if (isAllowed) {
                setDiscoveryDocument(discovery)
            }
        })
        return () => {
            isAllowed = false
        }
    }, [issuerOrDiscovery])
    return discoveryDocument
}


export function useLoadedAuthRequest(
    config: AuthRequestConfig,
    discovery: DiscoveryDocument | null,
    AuthRequestInstance: typeof AuthRequest
): AuthRequest | null {
    const [request, setRequest] = useState<AuthRequest | null>(null);
    const scopeString = useMemo(() => config.scopes?.join(','), [config.scopes]);
    const extraParamsString = useMemo(
        () => JSON.stringify(config.extraParams || {}),
        [config.extraParams]
    );
    useEffect(
        () => {
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
    return request;
}

type PromptMethod = (options?: AuthRequestPromptOptions) => Promise<AuthSessionResult>;

export function useAuthRequestResult(
    request: AuthRequest | null,
    discovery: DiscoveryDocument | null
): [AuthSessionResult | null, PromptMethod] {
    const [result, setResult] = useState<AuthSessionResult | null>(null);

    const promptAsync = useCallback(
        async () => {
            if (!discovery || !request) {
                throw new Error('Cannot prompt to authenticate until the request has finished loading.');
            }
            const result = await request?.promptAsync(discovery);
            setResult(result);
            return result;
        },
        // eslint-disable-next-line react-hooks/exhaustive-deps
        [request?.url, discovery?.authorization_endpoint]
    );

    return [result, promptAsync];
}

export function useAuthRequest(
    config: AuthRequestConfig,
    discovery: DiscoveryDocument | null
): [
        AuthRequest | null,
        AuthSessionResult | null,
    (options?: AuthRequestPromptOptions) => Promise<AuthSessionResult>
] {
    const request = useLoadedAuthRequest(config, discovery, AuthRequest);
    const [result, promptAsync] = useAuthRequestResult(request, discovery);
    return [request, result, promptAsync];
}

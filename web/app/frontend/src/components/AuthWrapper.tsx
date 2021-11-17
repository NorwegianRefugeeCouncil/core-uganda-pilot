import React, {FC, Fragment, useCallback, useEffect, useMemo, useState} from 'react';
import {axiosInstance} from "../utils/clients"
import {
    CodeChallengeMethod,
    exchangeCodeAsync,
    makeRedirectUri,
    ResponseType,
    TokenResponse,
    useAuthRequest,
    useAutoDiscovery,
} from 'expo-auth-session';
import {Button, Platform} from "react-native";
import Constants from "expo-constants";

export const AuthWrapper: FC = props => {
    const {children} = props
    const clientId = Constants.manifest?.extra?.client_id
    const useProxy = useMemo(() => Platform.select({web: false, default: false}), []);
    const redirectUri = useMemo(() => makeRedirectUri({scheme: 'nrccore'}), [])
    const discovery = useAutoDiscovery(Constants.manifest?.extra?.issuer);
    const [loggedIn, setLoggedIn] = useState(false)
    const [tokenResponse, setTokenResponse] = useState<TokenResponse>()

    const [request, response, promptAsync] = useAuthRequest(
        {
            clientId,
            usePKCE: true,
            responseType: ResponseType.Code,
            codeChallengeMethod: CodeChallengeMethod.S256,
            scopes: Constants.manifest?.extra?.scopes,
            redirectUri,
        },
        discovery
    );

    React.useEffect(() => {
        if (!discovery) {
            return;
        }
        if (!request?.codeVerifier) {
            return;
        }
        if (!response || response.type !== "success") {
            return;
        }

        const exchangeConfig = {
            code: response.params.code,
            clientId,
            redirectUri,
            extraParams: {
                "code_verifier": request?.codeVerifier,
            }
        }

        exchangeCodeAsync(exchangeConfig, discovery)
            .then(a => {
                setTokenResponse(a)
            })
            .catch((err) => {
                setTokenResponse(undefined)
            })

    }, [request?.codeVerifier, response, discovery]);

    useEffect(() => {
        if (!discovery) {
            return
        }
        if (tokenResponse?.shouldRefresh()) {
            const refreshConfig = {
                clientId: clientId,
                scopes: Constants.manifest?.extra?.scopes,
                extraParams: {}
            }
            tokenResponse?.refreshAsync(refreshConfig, discovery)
                .then(resp => {
                    setTokenResponse(resp)
                })
                .catch((err) => {
                    setTokenResponse(undefined)
                })
        }
    }, [tokenResponse?.shouldRefresh(), discovery])

    useEffect(() => {
        if (tokenResponse) {
            if (!loggedIn) {
                setLoggedIn(true)
            }
        } else {
            if (loggedIn) {
                setLoggedIn(false)
            }
        }
    }, [tokenResponse, loggedIn])


    useEffect(() => {
        const interceptor = axiosInstance.interceptors.request.use(value => {
            if (!tokenResponse?.accessToken) {
                return value
            }
            if (!value.headers) {
                value.headers = {}
            }
            value.headers["Authorization"] = `Bearer ${tokenResponse.accessToken}`
            return value
        })
        return () => {
            axiosInstance.interceptors.request.eject(interceptor)
        }
    }, [tokenResponse?.accessToken])

    const handleLogin = useCallback(() => {
        promptAsync({useProxy}).then(response => {
            console.log("PROMPT RESPONSE", response)
        }).catch((err) => {
            console.log("PROMPT ERROR", err)
        })
    }, [useProxy, promptAsync])

    if (!loggedIn) {
        return (
            <Button
                title={"Login"}
                disabled={!request}
                onPress={handleLogin}
            />
        )
    }
    return (
        <Fragment>
            {children}
        </Fragment>
    )
}

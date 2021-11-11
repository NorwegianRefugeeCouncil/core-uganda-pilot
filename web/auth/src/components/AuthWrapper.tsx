import React, {FC, Fragment, useCallback, useEffect, useState} from "react";
import {maybeCompleteAuthSession} from "../browser";
import {exchangeCodeAsync, TokenResponse} from "../tokenrequest";
import {AuthWrapperProps, CodeChallengeMethod, ResponseType} from "../types/types";
import {useAuthRequest, useDiscovery} from "../hooks";
import {axiosInstance} from "../client";

maybeCompleteAuthSession()

const AuthWrapper: FC<AuthWrapperProps> = ({
    children,
    scopes,
    clientId,
    issuer,
    redirectUriSuffix
}) => {
    const redirectUri = `${window.location.protocol}//${window.location.host}/${redirectUriSuffix}`;

    const discovery = useDiscovery(issuer)
    const [tokenResponse, setTokenResponse] = useState<TokenResponse>()
    const [isLoggedIn, setIsLoggedIn] = useState(false)

    const [request, response, promptAsync] = useAuthRequest(
        {
            clientId,
            usePKCE: true,
            responseType: ResponseType.Code,
            codeChallengeMethod: CodeChallengeMethod.S256,
            scopes,
            redirectUri
        },
        discovery
    );

    useEffect(() => {

        if (!discovery) {
            return
        }
        if (!request?.codeVerifier) {
            return
        }
        if (!response || response?.type !== "success") {
            return
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
                scopes: ["openid", "profile", "offline_access"],
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
            if (!isLoggedIn) {
                setIsLoggedIn(true)
            }
        } else {
            if (isLoggedIn) {
                setIsLoggedIn(false)
            }
        }
    }, [tokenResponse, isLoggedIn])

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
        promptAsync().then(response => {
            console.log("PROMPT RESPONSE", response)
        }).catch((err) => {
            console.log("PROMPT ERROR", err)
        })
    }, [discovery, request, promptAsync])

    if (!isLoggedIn) {
        return (
            <Fragment>
                <button onClick={handleLogin}>Login</button>
            </Fragment>
        )
    }

    return (
        <Fragment>
            {children}
        </Fragment>
    )

}

export default AuthWrapper;

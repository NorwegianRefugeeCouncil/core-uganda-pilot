import axios from "axios";
import React, {useCallback, useEffect, Fragment, useState} from "react";
import Browser from "../types/browser";
import exchangeCodeAsync from "../utils/exchangeCodeAsync";
import useDiscovery from "../hooks/useDiscovery";
import useAuthRequest from "../hooks/useAuthRequest";
import {AuthWrapperProps, CodeChallengeMethod, ResponseType} from "../types/types";
import {TokenResponse} from "../types/response";

const AuthWrapper: React.FC<AuthWrapperProps> = (
    {
        children,
        scopes = [],
        clientId,
        axiosInstance = axios.create(),
        issuer,
        redirectUri,
        customLoginComponent,
        handleLoginErr = console.log,
    }
) => {
    const browser = new Browser();
    browser.maybeCompleteAuthSession();

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
        discovery,
        browser
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
            .then(resp => {
                setTokenResponse(resp)
            })
            .catch((err) => {
                console.log('Code Exchange Error', err)
                setTokenResponse(undefined)
            })

    }, [request?.codeVerifier, response, discovery]);

    useEffect(() => {
        if (!discovery) {
            return
        }
        if (tokenResponse?.shouldRefresh()) {
            const refreshConfig = {
                clientId,
                scopes,
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
        promptAsync().catch((err) => {
            handleLoginErr(err);
        })
    }, [discovery, request, promptAsync])

    if (!isLoggedIn) {
        return (
            <Fragment>
                {customLoginComponent
                    ?
                    customLoginComponent({login: handleLogin})
                    :
                    <button onClick={handleLogin} role={"button"}>Login</button>
                }
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

import React from "react";
import {maybeCompleteAuthSession} from "../browser";
import {exchangeCodeAsync, TokenResponse} from "../tokenrequest";
import {AuthWrapperProps, CodeChallengeMethod, ResponseType} from "../types/types";
import {useAuthRequest, useDiscovery} from "../hooks";
import axios from "axios";

maybeCompleteAuthSession()

// TODO: https://betterprogramming.pub/building-secure-login-flow-with-oauth-2-openid-in-react-apps-ce6e8e29630a

const AuthWrapper: React.FC<AuthWrapperProps> = (
    {
        children,
        scopes = [],
        clientId,
        axiosInstance = axios.create(),
        issuer,
        redirectUriSuffix='/',
        customLoginComponent,
        handleLoginErr = console.log,
    }
) => {

    const redirectUri = `${window.location.protocol}//${window.location.host}/${redirectUriSuffix}`;

    const discovery = useDiscovery(issuer)
    const [tokenResponse, setTokenResponse] = React.useState<TokenResponse>()
    const [isLoggedIn, setIsLoggedIn] = React.useState(false)

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

    React.useEffect(() => {

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

    React.useEffect(() => {
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

    React.useEffect(() => {
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

    React.useEffect(() => {
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

    const handleLogin = React.useCallback(() => {
        promptAsync().catch((err) => {
            handleLoginErr(err);
        })
    }, [discovery, request, promptAsync])

    if (!isLoggedIn) {
        return (
            <React.Fragment>
                {customLoginComponent
                    ?
                    customLoginComponent({login: handleLogin})
                    :
                    <button onClick={handleLogin}>Login</button>
                }
            </React.Fragment>
        )
    }

    return (
        <React.Fragment>
            {children}
        </React.Fragment>
    )

}

export default AuthWrapper;

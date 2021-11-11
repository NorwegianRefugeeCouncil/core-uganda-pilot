import React, {FC, Fragment, useCallback, useEffect, useMemo, useState} from 'react';
import {Provider as PaperProvider} from 'react-native-paper';
import theme from './src/constants/theme';
import Router from './src/components/Router';
import {axiosInstance} from "./src/utils/clients"
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
import * as WebBrowser from 'expo-web-browser';
import Constants from "expo-constants";

WebBrowser.maybeCompleteAuthSession();


export const AuthWrapper: FC = props => {
    const {children} = props
    const clientId = `${Constants.manifest?.extra?.client_id}`
    const useProxy = useMemo(() => Platform.select({web: false, default: false}), []);
    const redirectUri = useMemo(() => makeRedirectUri({scheme: 'nrccore'}), [])
    const discovery = useAutoDiscovery(`${Constants.manifest?.extra?.issuer}`);
    const [loggedIn, setLoggedIn] = useState(false)

    const [request, response, promptAsync] = useAuthRequest(
        {
            clientId,
            usePKCE: true,
            responseType: ResponseType.Code,
            codeChallengeMethod: CodeChallengeMethod.S256,
            scopes: ['openid', 'profile', 'offline_access'],
            redirectUri,
        },
        discovery
    );

    const [tokenResponse, setTokenResponse] = useState<TokenResponse>()

    React.useEffect(() => {

        if (!discovery) {
            return
        }
        if (!request?.codeVerifier) {
            return
        }
        if (!response || response.type !== "success") {
            return
        }

        console.log('RESPONSE', response)

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
                console.log("EXCHANGE SUCCESS", a)
                setTokenResponse(a)
            })
            .catch((err) => {
                console.log("EXCHANGE ERROR", err)
                setTokenResponse(undefined)
            })

    }, [request?.codeVerifier, response, discovery]);

    useEffect(() => {
        if (!discovery) {
            return
        }
        if (tokenResponse?.shouldRefresh()) {
            console.log("REFRESHING TOKEN")
            const refreshConfig = {
                clientId: clientId,
                scopes: ["openid", "profile", "offline_access"],
                extraParams: {}
            }
            tokenResponse?.refreshAsync(refreshConfig, discovery)
                .then(resp => {
                    console.log("TOKEN REFRESH SUCCESS", resp)
                    setTokenResponse(resp)
                })
                .catch((err) => {
                    console.log("TOKEN REFRESH ERROR", err)
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
        console.log("SETTING UP INTERCEPTOR")
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
        return <PaperProvider theme={theme}>
            <Button
                title={"Login"}
                disabled={!request}
                onPress={handleLogin}
            />
        </PaperProvider>
    }
    return <Fragment>
        {children}
    </Fragment>
}

export default function App() {
    return (
        <PaperProvider theme={theme}>
            <AuthWrapper>
                <Router/>
            </AuthWrapper>
        </PaperProvider>
    );
}



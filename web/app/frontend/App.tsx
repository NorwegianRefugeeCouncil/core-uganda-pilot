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
    useAuthRequest,
    useAutoDiscovery,
} from 'expo-auth-session';
import {Button, Platform} from "react-native";
import * as WebBrowser from 'expo-web-browser';

WebBrowser.maybeCompleteAuthSession();


export const AuthWrapper: FC = props => {
    const {children} = props
    const clientId = 'react-native' //TODO
    const useProxy = useMemo(() => Platform.select({web: false, default: false}), []);
    const redirectUri = useMemo(() => makeRedirectUri({scheme: 'nrccore'}), [])
    const discovery = useAutoDiscovery('https://oidc.dev:8443');
    const [loggedIn, setLoggedIn] = useState(false)

    const [request, response, promptAsync] = useAuthRequest(
        {
            clientId,
            usePKCE: true,
            responseType: ResponseType.Code,
            codeChallengeMethod: CodeChallengeMethod.S256,
            scopes: ['openid', 'profile'],
            redirectUri,
        },
        discovery
    );

    const [accessToken, setAccessToken] = useState("")
    const [idToken, setIdToken] = useState("")
    const [refreshToken, setRefreshToken] = useState("")
    const [expiresIn, setExpiresIn] = useState<number | undefined>(undefined)
    const [issuedAt, setIssuedAt] = useState<number | undefined>(undefined)

    React.useEffect(() => {
        console.log('RESPONSE', response)

        if (!discovery) {
            return
        }
        if (!request || !request.codeVerifier) {
            return
        }
        if (!response || response.type !== "success") {
            return
        }

        exchangeCodeAsync({
            code: response.params.code,
            clientId,
            redirectUri,
            extraParams: {
                "code_verifier": request?.codeVerifier,
            }
        }, discovery)
            .then(a => {
                console.log("EXCHANGE SUCCESS", a)
                setIdToken(a.idToken ? a.idToken : "")
                setAccessToken(a.idToken ? a.accessToken : "")
                setRefreshToken(a.refreshToken ? a.refreshToken : "")
                setExpiresIn(a.expiresIn)
                setLoggedIn(true)
                setIssuedAt(a.issuedAt)
            })
            .catch((err) => {
                console.log("EXCHANGE ERROR", err)
            })
    }, [request, response, discovery]);

    const handleLogin = useCallback(() => {
        promptAsync({useProxy,}).then(response => {
            console.log("PROMPT RESPONSE", response)
        }).catch((err) => {
            console.log("PROMPT ERROR", err)
        })
    }, [useProxy, promptAsync])

    useEffect(() => {
        console.log("SETTING UP INTERCEPTOR")
        const interceptor = axiosInstance.interceptors.request.use(value => {
            console.log("INTERCEPTED VALUE")
            if (!value) {
                return
            }
            if (!accessToken) {
                return
            }
            if (!value.headers) {
                value.headers = {}
            }
            value.headers["Authorization"] = `Bearer ${accessToken}`
            console.log("NEW REQUEST", value)
            return value
        }, error => {
            console.error(error)
            return error
        })
        return () => {
            console.log("EJECTING INTERCEPTOR")
            axiosInstance.interceptors.request.eject(interceptor)
        }
    }, [accessToken])

    if (!loggedIn) {
        return <PaperProvider theme={theme}>
            <Button
                title={"Login"}
                disabled={!request}
                onPress={handleLogin}
            />
        </PaperProvider>
    }
    if (!accessToken) {
        return <Fragment/>
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



import React from 'react';
import {Provider as PaperProvider} from 'react-native-paper';
import theme from './src/constants/theme';
import Router from './src/components/Router';
import {
    CodeChallengeMethod,
    exchangeCodeAsync,
    makeRedirectUri,
    ResponseType,
    useAuthRequest,
    useAutoDiscovery
} from 'expo-auth-session';
import {Button, Platform} from "react-native";
import * as WebBrowser from 'expo-web-browser';

WebBrowser.maybeCompleteAuthSession();

export default function App() {

    // const [loggedIn, setLoggedIn] = React.useState(false);
    const useProxy = Platform.select({web: false, default: true});
    const discovery = useAutoDiscovery('http://localhost:4444');

    const [request, response, promptAsync] = useAuthRequest(
        {
            clientId: 'react-native',
            usePKCE: true,
            responseType: ResponseType.Code,
            codeChallengeMethod: CodeChallengeMethod.S256,
            codeChallenge:
            scopes: ['openid', 'profile'],
            // For usage in managed apps using the proxy
            redirectUri: makeRedirectUri({
                scheme: 'nrccore',

                // path: '/callback'
            }),


        },
        discovery
    );

    React.useEffect(() => {
        console.log('RESPONSE', response)
        if (response && response.type === 'success') {
            if (discovery==null) {
                return
            }
            const token = response.params.access_token;
            exchangeCodeAsync({
                code: response.params.code,
                clientId: 'react-native',
                redirectUri: makeRedirectUri({
                    scheme: 'nrccore',
                })
            }, discovery)
            console.log('TOKEN', token)
        }
    }, [response, discovery]);


    return (
        <PaperProvider theme={theme}>
            <Router/>
            <Button
                title={"Login"}
                disabled={!request}
                onPress={() => {
                    promptAsync({useProxy});
                }}
            />
        </PaperProvider>
    );
}



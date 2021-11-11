import React, {FC, Fragment, useState} from "react";
import {maybeCompleteAuthSession} from "core-auth/lib/browser";

// type SessionWrapperProps = {}

maybeCompleteAuthSession()

// type Metadata = {
//     authorization_endpoint: string
//     token_endpoint: string
// }
//
// type ExchangeResponse = {
//     access_token: string
//     id_token: string
//     refresh_token: string
//     expires_in: number
//     scope: string
//     token_type: string
// }

const AuthWrapper: FC = (props) => {
    // const isLoggedIn = true;
    // let tokenResponse;
    // let discovery;
    // let request;
    // let response;

    // const discovery = useDiscovery(`${process.env.REACT_APP_ISSUER}`)
    const clientId = `${process.env.REACT_APP_CLIENT_ID}`
    const redirectUri = `${global.window.location.protocol}//${global.window.location.host}/app`;
    // const [accessToken, setAccessToken] = useState("")
    // const [tokenResponse, setTokenResponse] = useState<TokenResponse>()
    const [isLoggedIn, setIsLoggedIn] = useState(false)

    // const [request, response, promptAsync] = useAuthRequest(
    //     {
    //         clientId,
    //         usePKCE: true,
    //         responseType: ResponseType.Code,
    //         codeChallengeMethod: CodeChallengeMethod.S256,
    //         scopes: ['openid', 'profile', 'offline_access'],
    //         redirectUri
    //     },
    //     discovery
    // );

    //
    // useEffect(() => {
    //
    //     if (!discovery) {
    //         return
    //     }
    //     if (!request?.codeVerifier) {
    //         return
    //     }
    //     if (!response || response?.type !== "success") {
    //         return
    //     }
    //
    //     // console.log('RESPONSE', response)
    //
    //     const exchangeConfig = {
    //         code: response.params.code,
    //         clientId,
    //         redirectUri,
    //         extraParams: {
    //             "code_verifier": request?.codeVerifier,
    //         }
    //     }
    //
    //     exchangeCodeAsync(exchangeConfig, discovery)
    //         .then(a => {
    //             // console.log("EXCHANGE SUCCESS", a)
    //             // setTokenResponse(a)
    //         })
    //         .catch((err) => {
    //             // console.log("EXCHANGE ERROR", err)
    //             // setTokenResponse(undefined)
    //         })
    //
    // }, [request?.codeVerifier, response, discovery]);
    //
    // useEffect(() => {
    //     if (!discovery) {
    //         return
    //     }
    //     if (tokenResponse?.shouldRefresh()) {
    //         console.log("REFRESHING TOKEN")
    //         const refreshConfig = {
    //             clientId: clientId,
    //             scopes: ["openid", "profile", "offline_access"],
    //             extraParams: {}
    //         }
    //         tokenResponse?.refreshAsync(refreshConfig, discovery)
    //             .then(resp => {
    //                 // console.log("TOKEN REFRESH SUCCESS", resp)
    //                 // setTokenResponse(resp)
    //             })
    //             .catch((err) => {
    //                 // console.log("TOKEN REFRESH ERROR", err)
    //                 // setTokenResponse(undefined)
    //             })
    //     }
    // }, [tokenResponse?.shouldRefresh(), discovery])
    //
    // useEffect(() => {
    //     if (tokenResponse) {
    //         if (!isLoggedIn) {
    //             // setIsLoggedIn(true)
    //         }
    //     } else {
    //         if (isLoggedIn) {
    //             // setIsLoggedIn(false)
    //         }
    //     }
    // }, [tokenResponse, isLoggedIn])
    //
    // useEffect(() => {
    //     // console.log("SETTING UP INTERCEPTOR")
    //     const interceptor = axiosInstance.interceptors.request.use(value => {
    //         if (!tokenResponse?.accessToken) {
    //             return value
    //         }
    //         if (!value.headers) {
    //             value.headers = {}
    //         }
    //         value.headers["Authorization"] = `Bearer ${tokenResponse.accessToken}`
    //
    //         return value
    //     })
    //     return () => {
    //         axiosInstance.interceptors.request.eject(interceptor)
    //     }
    // }, [tokenResponse?.accessToken])
    //
    // const handleLogin = useCallback(() => {
    //     // promptAsync().then(response => {
    //     //     console.log("PROMPT RESPONSE", response)
    //     // }).catch((err) => {
    //     //     console.log("PROMPT ERROR", err)
    //     // })
    // }, [1])


    // useEffect(() => {
    //     if (accessToken) {
    //         setIsLoggedIn(true)
    //     } else {
    //         setIsLoggedIn(false)
    //     }
    //     const int = axiosInstance.interceptors.request.use(value => {
    //         console.log("INTERCEPTED REQUEST")
    //         if (!accessToken) {
    //             return value
    //         }
    //         if (!value.headers) {
    //             value.headers = {}
    //         }
    //         value.headers["Authorization"] = `Bearer ${accessToken}`
    //         return value
    //     })
    //     return () => axiosInstance.interceptors.request.eject(int)
    // }, [accessToken])

    // const {children} = props
    // if (!isLoggedIn) {
    //     return (
    //         <Fragment>
    //
    //         </Fragment>
    //     )
    // }

    return (
        <Fragment>
            {!isLoggedIn && <button>Login</button>}
            {isLoggedIn && props.children}
        </Fragment>
    )

}

export default AuthWrapper;

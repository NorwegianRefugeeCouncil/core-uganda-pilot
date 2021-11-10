import {FC, Fragment, useCallback, useEffect, useState} from "react";
import axios, {AxiosRequestConfig} from "axios";
import getPkce from "oauth-pkce"
import {axiosInstance} from "../app/client";

type SessionWrapperProps = {}

type Metadata = {
    authorization_endpoint: string
    token_endpoint: string
}

type ExchangeResponse = {
    access_token: string
    id_token: string
    refresh_token: string
    expires_in: number
    scope: string
    token_type: string
}

export const AuthWrapper: FC<SessionWrapperProps> = props => {
    const {children} = props
    const [pkce, setPkce] = useState<{ challenge: string, verifier: string }>()
    const [metadata, setMetadata] = useState<Metadata>()
    const [state, setState] = useState("")
    const [isLoggedIn, setIsLoggedIn] = useState(false)
    const [accessToken, setAccessToken] = useState("")
    const [idToken, setIdToken] = useState("")
    const [refreshToken, setRefreshToken] = useState("")
    const [expiresIn, setExpiresIn] = useState(0)
    const [pending, setPending] = useState(false)

    useEffect(() => {
        getPkce(43, (error, {verifier, challenge}) => {
            if (error) {
                console.log("PKCE ERROR", error)
                return
            }
            setPkce({verifier, challenge})
        })
    }, [])

    useEffect(() => {
        const metadataEndpoint = `${process.env.REACT_APP_ISSUER}/.well-known/openid-configuration`;
        axios.get(metadataEndpoint)
            .then(value => value.data as Metadata)
            .then(value => setMetadata(value))
    }, [])

    useEffect(() => {
        const arr = new Uint8Array(32)
        setState(btoa(crypto.getRandomValues(arr).toString()))
    }, [])

    const doLogin = useCallback(() => {
        if (!metadata) {
            return
        }
        if (!pkce) {
            return
        }
        if (!state) {
            return
        }
        if (pending) {
            return
        }
        const redirectUri = `${window.location.protocol}//${window.location.host}/app`;
        const params = [
            {
                key: "scope",
                value: "openid profile email offline_access"
            }, {
                key: "state",
                value: state,
            }, {
                key: "redirect_uri",
                value: redirectUri,
            }, {
                key: "response_type",
                value: "code",
            }, {
                key: "client_id",
                value: `${process.env.REACT_APP_CLIENT_ID}`,
            }, {
                key: "code_challenge",
                value: `${pkce.challenge}`,
            }, {
                key: "code_challenge_method",
                value: `S256`,
            },
        ]
            .map(({key, value}) => ({key: key, value: encodeURIComponent(value)}))
            .map(({key, value}) => `${key}=${value}`)
            .join("&")

        const authURL = `${metadata.authorization_endpoint}?${params}`
        const tokenURL = `${metadata.token_endpoint}?${params}`

        const popupWidth = 380
        const popupHeight = 480
        const left = window.screen.width / 2 - popupWidth / 2
        const top = window.screen.height / 2 - popupHeight / 2
        const popup = window.open(
            authURL,
            "Core Login",
            "menubar=no,location=no,resizable=no,scrollbars=no,status=no, width=" +
            popupWidth +
            ", height=" +
            popupHeight +
            ", top=" +
            top +
            ", left=" +
            left
        )
        setPending(true)

        if (popup !== null) {

            popup.onunload = ev => {
                console.log(popup.location)
                const outParams = new URLSearchParams(popup.location.search)
                const code = outParams.get("code")
                if (!code) {
                    return
                }
                const params = new URLSearchParams()
                params.set("code", `${code}`)
                params.set("grant_type", "authorization_code")
                params.set("client_id", `${process.env.REACT_APP_CLIENT_ID}`)
                params.set("code_verifier", `${pkce.verifier}`)
                params.set("redirect_uri", `${redirectUri}`)
                const config: AxiosRequestConfig = {
                    headers: {"Content-Type": "application/x-www-form-urlencoded"},
                }
                axios.post<ExchangeResponse>(tokenURL, params, config).then((resp) => {
                    console.log("TOKEN EXCHANGE SUCCESS", resp.data)
                    setAccessToken(resp.data.access_token)
                    setIdToken(resp.data.id_token)
                    setExpiresIn(resp.data.expires_in)
                    setRefreshToken(resp.data.refresh_token)
                    setPending(false)
                }, err => {
                    console.log("TOKEN EXCHANGE ERROR", err)
                    setPending(false)
                })
            }
        }
    }, [metadata, pkce, state, pending])

    useEffect(() => {
        if (!expiresIn) {
            return
        }
        if (!metadata) {
            return
        }
        console.log("EXPIRES IN", expiresIn)
        const renewalIn = Math.abs(Math.round(expiresIn / 2))
        console.log("RENEWAL IN", renewalIn)
        const timeout = setTimeout(() => {

            const params = new URLSearchParams()
            params.set("refresh_token", `${refreshToken}`)
            params.set("grant_type", "refresh_token")
            params.set("client_id", `${process.env.REACT_APP_CLIENT_ID}`)
            const config: AxiosRequestConfig = {
                headers: {"Content-Type": "application/x-www-form-urlencoded"},
            }
            axios.post<ExchangeResponse>(metadata.token_endpoint, params, config).then(resp => {
                console.log("TOKEN RENEWAL SUCCESS", resp)
                if (resp.data.access_token) {
                    setAccessToken(resp.data.access_token)
                }
                if (resp.data.refresh_token) {
                    setRefreshToken(resp.data.refresh_token)
                }
                if (resp.data.id_token) {
                    setIdToken(resp.data.id_token)
                }
                setExpiresIn(resp.data.expires_in)
            }, err => {
                console.log("TOKEN RENEWAL ERROR", err)
            })

        }, renewalIn * 1000)
        return () => clearTimeout(timeout)
    }, [metadata, expiresIn, refreshToken])

    useEffect(() => {
        if (accessToken) {
            setIsLoggedIn(true)
        } else {
            setIsLoggedIn(false)
        }
        const int = axiosInstance.interceptors.request.use(value => {
            console.log("INTERCEPTED REQUEST")
            if (!accessToken) {
                return value
            }
            if (!value.headers) {
                value.headers = {}
            }
            value.headers["Authorization"] = `Bearer ${accessToken}`
            return value
        })
        return () => axios.interceptors.request.eject(int)
    }, [accessToken])

    if (!isLoggedIn) {
        return <Fragment>
            {isLoggedIn ? "logged in" : "not logged in"}
            {pending ? "pending" : ""}
            <button onClick={doLogin} disabled={pending}>Login</button>
        </Fragment>
    }

    return <Fragment>
        {children}
    </Fragment>

}

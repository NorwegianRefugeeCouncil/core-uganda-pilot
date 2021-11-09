import {FC, Fragment, useCallback, useEffect, useState} from "react";
import axios from "axios";
import getPkce from "oauth-pkce"

type SessionWrapperProps = {}

type Metadata = {
    authorization_endpoint: string
    token_endpoint: string
}

export const AuthWrapper: FC<SessionWrapperProps> = props => {
    const {children} = props
    const [pkce, setPkce] = useState<{ challenge: string, verifier: string }>()
    const [metadata, setMetadata] = useState<Metadata>()
    const [state, setState] = useState("")
    const [isLoggedIn, setIsLoggedIn] = useState(false)

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

        const params = [
            {
                key: "scope",
                value: "openid profile email"
            }, {
                key: "state",
                value: state,
            }, {
                key: "redirect_uri",
                value: `${window.location.protocol}//${window.location.host}`,
            }, {
                key: "response_type",
                value: "token",
            }, {
                key: "client_id",
                value: `${process.env.REACT_APP_CLIENT_ID}`,
            },
        ]
            .map(({key, value}) => ({key: key, value: encodeURIComponent(value)}))
            .join("&")

        const authURL = `${metadata.authorization_endpoint}?${params}`

        const popupWidth = 380
        const popupHeight = 480
        const left = window.screen.width / 2 - popupWidth / 2
        const top = window.screen.height / 2 - popupHeight / 2
        const win = window.open(
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
        if (win) {
            win.opener = window
        }
    }, [metadata, pkce, state])

    if (!isLoggedIn) {
        return <button onClick={doLogin}>Login</button>
    }

    return <Fragment>{children}</Fragment>

}

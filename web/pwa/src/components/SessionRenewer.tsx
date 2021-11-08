import {FC, Fragment, useCallback, useEffect, useState} from "react";
import {useApiClient} from "../app/hooks";
import {Session} from "../types/types";
import {useLocation} from "react-router-dom";

export const SessionRenewer: FC = props => {

    const apiClient = useApiClient()
    const [session, setSession] = useState<Session>()
    const [refresh, setRefresh] = useState(true)

    const location = useLocation()

    const refreshSession = useCallback(() => {
        if (apiClient && refresh) {
            apiClient.getSession().then(s => {
                setSession(s.response)
                setRefresh(false)
            })
        }
    }, [apiClient, refresh])

    useEffect(() => {
        refreshSession()
    }, [refreshSession])

    useEffect(() => {
        const params = new URLSearchParams(location.search)
        if (params.get("error") == "login_required") {
            window.parent.location.href = `${apiClient.address}/oidc/login?redirect_uri=${window.parent.location.href}`
        }
    }, [location.search])

    useEffect(() => {
        if (!session) {
            return
        }
        const timeout = setInterval(() => {
            const expiry = new Date(session.expiry)
            const now = new Date()
            const expiresInSeconds =Math.round((expiry.getTime() - now.getTime()) / 1000)
            const renewalInSeconds = expiresInSeconds - 60
            console.log("session renews in :", renewalInSeconds)
            if (renewalInSeconds < 0) {
                window.location.href = `${apiClient.address}/oidc/renew?redirect_uri=${window.location.href}`
            }
        }, 2000)
        return () => {
            clearInterval(timeout)
        }
    }, [session])


    return <Fragment/>

}

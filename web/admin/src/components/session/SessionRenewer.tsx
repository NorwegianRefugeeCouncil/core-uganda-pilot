import {FC, Fragment, useCallback, useEffect, useState} from "react";
import {useApiClient} from "../../app/hooks";
import {Session} from "../../client/client";

export const SessionRenewer: FC = props => {

    const apiClient = useApiClient()
    const [session, setSession] = useState<Session>()
    const [refresh, setRefresh] = useState(true)

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
        if (!session) {
            return
        }
        const timeout = setInterval(() => {
            const expiry = new Date(session.expiry)
            const now = new Date()
            const expiresInSeconds = (expiry.getTime() - now.getTime()) / 1000
            const renewalInSeconds = expiresInSeconds - 150
            console.log("session renews in :", renewalInSeconds)
            if (renewalInSeconds < 0) {
                setTimeout(() => {
                    window.location.href = `${apiClient.address}/oidc/renew?redirect_uri=${window.location.href}`
                }, 2000)
            }
        }, 2000)
        return () => {
            clearInterval(timeout)
        }
    }, [session])


    return <Fragment/>

}

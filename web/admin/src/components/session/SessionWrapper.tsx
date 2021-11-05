import {FC, Fragment, useCallback, useEffect, useState} from "react";
import {Session} from "../../client/client";
import {useApiClient} from "../../app/hooks";
import axios from "axios";

type SessionWrapperProps = {}

export const SessionWrapper: FC<SessionWrapperProps> = props => {

    const {children} = props

    const [session, setSession] = useState<Session>()
    const [isRefreshing, setIsRefreshing] = useState(false)
    const apiClient = useApiClient()

    const refresh = useCallback(() => {
        if (isRefreshing) {
            return
        }
        apiClient.getSession().then(session => {
            setSession(session.response)
            if (!session.response?.active) {
                setTimeout(() => {
                    window.location.href = `${apiClient.address}/oidc/login?redirect_uri=${window.location.href}`
                }, 2000)
            }
        }).finally(() => {
            setIsRefreshing(false)
        })
    }, [isRefreshing])

    useEffect(() => {
        axios.interceptors.response.use(
            value => value,
            error => {
                refresh()
                return error
            })
    }, [])

    useEffect(() => {
        refresh()
    }, [])

    if (!session?.active) {
        return <div>Authenticating...</div>
    }

    return <Fragment>{children}</Fragment>

}

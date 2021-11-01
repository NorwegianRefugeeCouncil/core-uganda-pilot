import {useContext, useEffect, useMemo, useState} from "react";
import {AuthContext} from "oidc-react";
import {client, IdentityProvider, Organization} from "../client/client";
import {AuthContextProps} from "oidc-react/build/src/AuthContextInterface";
import {useParams} from "react-router-dom";
import {OrganizationRoute} from "../components/organization/OrganizationPortal";

export function useAuth(): AuthContextProps | undefined {
    return useContext(AuthContext)
}

export function useIsAuthenticated(): boolean {
    let authCtx = useAuth();
    return !!(authCtx?.userData)
}

export function useApiClient(): client {
    const authCtx = useAuth()
    return useMemo(() => {
        return new client({idToken: authCtx?.userData?.id_token})
    }, [authCtx?.userData?.id_token])
}

export function usePathOrganization(): Organization | undefined {
    const apiClient = useApiClient()
    const {organizationId} = useParams<OrganizationRoute>()
    const [organization, setOrganization] = useState<Organization>()
    useEffect(() => {
        if (!organizationId) {
            return
        }
        apiClient.getOrganization({id: organizationId}).then(resp => {
            if (resp.response) {
                setOrganization(resp.response)
            }
        })
    }, [apiClient, organizationId])
    return organization
}

export function useIdentityProviders(organizationId: string): IdentityProvider[] {
    const apiClient = useApiClient()
    const [idps, setIdps] = useState<IdentityProvider[]>([])
    useEffect(() => {
        if (!organizationId) {
            return
        }
        apiClient.listIdentityProviders({organizationId}).then(resp => {
            if (resp.response) {
                setIdps(resp.response ? resp.response.items : [])
            }
        })
    }, [apiClient, organizationId])
    return idps
}

import {useCallback, useContext, useEffect, useMemo, useState} from "react";
import {AuthContext} from "oidc-react";
import {client, IdentityProvider, Organization} from "../client/client";
import {AuthContextProps} from "oidc-react/build/src/AuthContextInterface";
import {useParams} from "react-router-dom";
import {OrganizationRoute} from "../components/organization/OrganizationPortal";
import classNames from "classnames";
import {UseFormReturn} from "react-hook-form";

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

export function useFormValidation<T extends object = { [key: string]: any }>(isNew: boolean, form: UseFormReturn<T>) {

    const {formState: {dirtyFields, errors, isSubmitted, touchedFields}} = form

    const fieldClasses = useCallback((field: keyof T) => {
        let cls = classNames("form-control form-control-darkula")
        const touched = touchedFields.hasOwnProperty(field)
        const dirty = dirtyFields.hasOwnProperty(field)
        const hasError = errors.hasOwnProperty(field)
        if (isSubmitted || (isNew && touched) || (isNew && dirty)) {
            if (hasError) {
                return classNames(cls, "is-invalid")
            }
        }
        return cls
    }, [dirtyFields, errors, isNew, isSubmitted, touchedFields])

    const fieldErrors = useCallback((field: keyof T) => {
        const touched = touchedFields.hasOwnProperty(field)
        const dirty = dirtyFields.hasOwnProperty(field)
        const hasError = errors.hasOwnProperty(field)
        const err = (errors as any)[field]
        return hasError && (isSubmitted || (isNew && touched) || (isNew && dirty)) &&
            <div className={"invalid-feedback"}>
                {err?.type === "required" && <span>This field is required</span>}
                {err?.type === "pattern" && <span>Invalid value</span>}
            </div>
    }, [touchedFields, dirtyFields, errors, isNew])

    return {fieldErrors, fieldClasses}

}

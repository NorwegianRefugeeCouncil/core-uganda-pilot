import React, {FC, useCallback, useEffect} from "react";
import {useAppDispatch, useAppSelector} from "../../app/hooks";
import {fetchIdentityProviders, identityProviderGlobalSelectors} from "../../reducers/identityproviders";
import {IdentityProvider} from "../../types/types";
import {Link, useRouteMatch} from "react-router-dom";

export type IdentityProvidersProps = {
    organizationId: string
}


export const IdentityProviders: FC<IdentityProvidersProps> = props => {

    const match = useRouteMatch()

    const dispatch = useAppDispatch()
    useEffect(() => {
        dispatch(fetchIdentityProviders({organizationId: props.organizationId}))
    }, [dispatch, props.organizationId])

    const idps = useAppSelector(state => {
        return identityProviderGlobalSelectors.selectForOrganization(state, props.organizationId)
    })

    const renderIdp = useCallback((idp: IdentityProvider) => {
        return <li key={idp.id} className={"list-group-item"}>
            <Link to={`${match.url}/${idp.id}`}>{idp.kind} - {idp.domain}</Link>
        </li>
    }, [match.url])

    const renderIdps = useCallback(() => {
        if (idps.length) {
            return <ul className={"list-group"}>
                {idps.map(renderIdp)}
            </ul>
        }
        return <h6 className={"text-muted"}>
            No identity provider currently registered for organization
        </h6>

    }, [renderIdp, idps])

    return (
        <div>
            <div className={"mb-2"}>
                <Link to={`${match.url}/new`} className={"btn btn-primary"}>Add Identity Provider</Link>
            </div>
            {renderIdps()}
        </div>
    )

}

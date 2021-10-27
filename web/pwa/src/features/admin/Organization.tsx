import React, {FC, Fragment, useEffect} from "react";
import {useAppSelector} from "../../app/hooks";
import {fetchOrganizations, organizationGlobalSelectors} from "../../reducers/organizations";
import {Link, useParams, useRouteMatch, Switch, Route} from "react-router-dom"
import {IdentityProviders} from "./IdentityProviders";
import {useDispatch} from "react-redux";
import {IdentityProviderEditorContainer} from "./IdentityProviderEditor";

export const Organization: FC = props => {

    const match = useRouteMatch()
    const params = useParams<{ organizationId: string }>()
    const dispatch = useDispatch()

    useEffect(() => {
        dispatch(fetchOrganizations())
    }, [dispatch])

    const organization = useAppSelector(state => {
        return organizationGlobalSelectors.selectById(state, params.organizationId)
    })

    if (!organization) {
        return <Fragment/>
    }

    return (
        <div className={"d-flex flex-row w-100 h-100 flex-grow-1"}>
            <ul className={"list-group h-100 border-end pt-2 rounded-0"}>
                <li className={"list-group-item border-0 rounded-0"}>
                    <Link
                        className={"btn btn-light"}
                        to={`${match.url}/identityproviders`}>Identity Providers</Link>
                </li>
            </ul>
            <div className={"flex-grow-1 bg-dark text-white"}>

                <h6 className={"border-bottom w-100 p-3"}>
                    <code className={"fw-bold"}>{organization.key}</code> {organization.name}
                </h6>

                <div className={"d-flex flex-column flex-grow-1 bg-dark p-3"}>
                    <Switch>
                        <Route exact path={`${match.url}/identityproviders/new`}>
                            <IdentityProviderEditorContainer organization={organization}/>
                        </Route>
                        <Route exact path={`${match.url}/identityproviders/:identityProviderId`}>
                            <IdentityProviderEditorContainer organization={organization}/>
                        </Route>
                        <Route exact path={`${match.url}/identityproviders`}>
                            <IdentityProviders organizationId={organization.id}/>
                        </Route>
                    </Switch>
                </div>

            </div>
        </div>
    )

}

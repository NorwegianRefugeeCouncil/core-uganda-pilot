import {FC, Fragment} from "react";
import {match, Route, Switch, useRouteMatch} from "react-router-dom";
import {usePathOrganization} from "../../hooks/hooks";
import {OrganizationOverview} from "./OrganizationOverview";
import {OrganizationSideBar} from "./OrganizationSideBar";
import {IdentityProviders} from "./identityproviders/IdentityProviders";
import {Organization} from "../../types/types";
import {IdentityProviderEditor} from "./identityproviders/IdentityProviderEditor";

export type OrganizationRoute = {
    organizationId: string
}


export const OrganizationPortal: FC = props => {

    const match = useRouteMatch()
    const organization = usePathOrganization()

    if (!organization) {
        return <Fragment/>
    }

    return (
        <div className={"flex-grow-1 d-flex flex-column"}>
            <div className={"py-2 ps-4 bg-darkula text-white"}>
                <h5 className={"p-0 m-2"}>{organization.name}</h5>
            </div>
            <div className={"d-flex flex-row flex-grow-1 mt-4 px-4"}>
                <div className={""}>
                    <OrganizationSideBar organization={organization}/>
                </div>
                <div className={"flex-grow-1 ps-4 pe-2"}>
                    <Switch>
                        {addIdentityProvidersRoute(match, organization)}
                        {identityProviderRoute(match, organization)}
                        {identityProvidersRoute(match, organization)}
                        {overviewRoute(match, organization)}
                    </Switch>
                </div>
            </div>
        </div>
    )
}


function identityProvidersRoute(m: match<{}>, organization: Organization) {
    return <Route path={`${m.path}/identity-providers`}
                  render={() => (<IdentityProviders organization={organization}/>)}/>;
}

function addIdentityProvidersRoute(m: match<{}>, organization: Organization) {
    return <Route path={`${m.path}/identity-providers/add`}
                  render={() => (<IdentityProviderEditor organization={organization}/>)}/>;
}

function identityProviderRoute(m: match<{}>, organization: Organization) {
    return <Route path={`${m.path}/identity-providers/:id`}
                  render={(p) => (
                      <IdentityProviderEditor id={p.match.params["id"]} organization={organization}/>
                  )}/>;
}

function overviewRoute(m: match<{}>, organization: Organization) {
    return <Route
        exact
        path={`${m.path}`}
        render={() => <OrganizationOverview organization={organization}/>}/>;
}

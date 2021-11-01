import {FC, Fragment} from "react";
import {Organization} from "../../../client/client";
import {useIdentityProviders} from "../../../app/hooks";
import {Link, useRouteMatch} from "react-router-dom";

type Props = {
    organization: Organization
}

export const IdentityProviders: FC<Props> = props => {

    const idps = useIdentityProviders(props.organization.id)
    const match = useRouteMatch()

    return (
        <div>
            <div className={"border-bottom border-secondary pb-3 my-2 d-flex flex-row text-light justify-content-center"}>
                <span className={"flex-grow-1 fs-5"}>Identity Providers</span>
                <Link className={"btn btn-darkula btn-sm"} to={`${match.path}/add`}>Add Identity Provider</Link>
            </div>

            <div className={"list-group list-group-darkula"}>
                {idps.map(idp => (
                    <Link className={"list-group-item list-group-item-action"}
                          to={`${match.url}/${idp.id}`}>{idp.name}</Link>
                ))}
                {idps.length === 0
                    ? <div className={"disabled list-group-item"}>No Identity Provider</div>
                    : <Fragment/>}
            </div>
        </div>
    )
}

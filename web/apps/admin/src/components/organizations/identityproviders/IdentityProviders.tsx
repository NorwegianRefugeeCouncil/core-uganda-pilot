import {FC, Fragment} from "react";
import {Organization} from "../../../types/types";
import {useIdentityProviders} from "../../../hooks/hooks";
import {Link, useRouteMatch} from "react-router-dom";
import {SectionTitle} from "../../sectiontitle/SectionTitle";

type Props = {
    organization: Organization
}

export const IdentityProviders: FC<Props> = props => {

    const idps = useIdentityProviders(props.organization.id)
    const match = useRouteMatch()

    return (
        <div>

            <SectionTitle className={"text-light"} title={"Identity Providers"}>
                <Link className={"btn btn-success btn-sm"} to={`${match.path}/add`}>Add Identity Provider</Link>
            </SectionTitle>

            <div className={"list-group list-group-darkula"}>
                {idps.map(idp => (
                    <Link className={"list-group-item list-group-item-action"}
                          to={`${match.url}/${idp.id}`}>
                        {idp.name} <span className={"badge bg-dark font-monospace"}>{idp.emailDomain}</span>
                    </Link>
                ))}
                {idps.length === 0
                    ? <div className={"disabled list-group-item"}>No Identity Provider</div>
                    : <Fragment/>}
            </div>
        </div>
    )
}

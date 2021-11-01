import {FC} from "react";
import {NavLink, useRouteMatch} from "react-router-dom"
import {Organization} from "../../client/client";
import styled from "styled-components"

type Props = {
    organization: Organization
}

export const OrganizationSideBar: FC<Props> = props => {
    const match = useRouteMatch()
    return (
        <div className={"list-group list-group-darkula"} style={{width: "15rem"}}>
            <NavLink
                activeClassName={"active"}
                exact={true}
                className={"list-group-item list-group-item-action"}
                to={`${match.url}`}>
                Overview
            </NavLink>

            <NavLink activeClassName={"active"}
                     className={"list-group-item list-group-item-action"}
                     to={`${match.url}/identity-providers`}>
                Identity Providers
            </NavLink>

        </div>
    )
}

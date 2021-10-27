import React, {FC} from "react";
import {Redirect, Route, Switch, useRouteMatch} from "react-router-dom";
import {Organizations} from "./Organizations";
import {AddOrganization} from "./AddOrganization";
import {Organization} from "./Organization";

export const AdminPanel: FC = props => {

    const match = useRouteMatch()

    console.log(match)

    return <div className={"d-flex flex-row w-100 h-100 flex-grow-1"}>


        <Switch>
            <Route path={`${match.url}/organizations`}
                   exact
                   component={Organizations}/>
            <Route path={`${match.url}/organizations/add`}
                   exact
                   component={AddOrganization}/>
            <Route path={`${match.url}/organizations/:organizationId`}
                   component={Organization}/>

            <Redirect to={`${match.url}/organizations`}/>

        </Switch>

    </div>

}

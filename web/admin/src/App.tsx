import React from 'react';
import './App.scss';
import 'bootstrap/dist/css/bootstrap.css';
import 'bootstrap-icons/font/bootstrap-icons.css';
import 'bootstrap/dist/js/bootstrap.bundle.min';
import {BrowserRouter, Route, Switch} from "react-router-dom";
import {Organizations} from "./components/organizations/Organizations";
import {OrganizationPortal} from "./components/organizations/OrganizationPortal";
import {OrganizationEditor} from "./components/organizations/OrganizationEditor";
import {Clients} from "./components/clients/Clients";
import {ClientEditor} from "./components/clients/ClientEditor";
import {SessionWrapper} from "./components/session/SessionWrapper";
import {SessionRenewer} from "./components/session/SessionRenewer";
import {NavBar} from "./components/navbar/NavBar";


function AuthenticatedApp() {
    return (
        <div className={"d-flex flex-column vh-100 vw-100 bg-dark"}>
            <NavBar/>
            <Switch>
                <Route path={"/organizations/add"} component={OrganizationEditor}/>
                <Route path={"/organizations/:organizationId"} component={OrganizationPortal}/>
                <Route path={"/organizations"} component={Organizations}/>
                <Route path={"/clients/add"} component={ClientEditor}/>
                <Route path={"/clients/:clientId"} component={ClientEditor}/>
                <Route path={"/clients"} component={Clients}/>
            </Switch>
        </div>
    )
}

function App() {
    return (
        <BrowserRouter>
            <Switch>
                <Route path={"/session-renew"} exact render={props => {
                    return <SessionRenewer/>
                }}/>
                <Route path={""} render={props => {
                    return (
                        <SessionWrapper>
                            <iframe
                                title={"login"}
                                style={{position: "absolute", top: 0, left: 0, visibility: "hidden"}}
                                src={`${window.location.protocol}//${window.location.host}/session-renew`}>
                            </iframe>
                            <AuthenticatedApp/>
                        </SessionWrapper>
                    )
                }}/>
            </Switch>
        </BrowserRouter>
    );
}

export default App;

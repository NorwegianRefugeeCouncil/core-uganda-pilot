import React, {Fragment} from 'react';
import './App.scss';
import 'bootstrap/dist/css/bootstrap.css';
import 'bootstrap-icons/font/bootstrap-icons.css';
import 'bootstrap/dist/js/bootstrap.bundle.min';
import {AuthProvider} from "oidc-react"
import {NavBar} from "./components/navbar/NavBar";
import {BrowserRouter, Switch} from "react-router-dom";
import {Organizations} from "./components/organizations/Organizations";
import {ProtectedRoute} from "./components/guard/ProtectedRoute";
import {OrganizationPortal} from "./components/organizations/OrganizationPortal";
import {OrganizationEditor} from "./components/organizations/OrganizationEditor";
import {Clients} from "./components/clients/Clients";
import {ClientEditor} from "./components/clients/ClientEditor";

const oidcConfig = {
    authority: 'https://dev-53701279.okta.com',
    clientId: '0oa2gb25jhuFDIkhd5d7',
    redirectUri: 'http://localhost:3001'
};

function App() {
    return (
        <Fragment>
            <AuthProvider
                scope={"openid profile email"}
                autoSignIn={true}
                {...oidcConfig} >
                <div className={"d-flex flex-column vh-100 vw-100 bg-dark"}>
                    <BrowserRouter>
                        <NavBar/>
                        <Switch>
                            <ProtectedRoute path={"/organizations/add"} component={OrganizationEditor}/>
                            <ProtectedRoute path={"/organizations/:organizationId"} component={OrganizationPortal}/>
                            <ProtectedRoute path={"/organizations"} component={Organizations}/>
                            <ProtectedRoute path={"/clients/add"} component={ClientEditor}/>
                            <ProtectedRoute path={"/clients/:clientId"} component={ClientEditor}/>
                            <ProtectedRoute path={"/clients"} component={Clients}/>
                        </Switch>
                    </BrowserRouter>
                </div>
            </AuthProvider>
        </Fragment>
    );
}

export default App;

import React, {Fragment} from 'react';
import './App.scss';
import 'bootstrap/dist/css/bootstrap.css';
import 'bootstrap-icons/font/bootstrap-icons.css';
import 'bootstrap/dist/js/bootstrap.bundle.min';
import {NavBar} from "./components/navbar/NavBar";
import {BrowserRouter, Route, Switch} from "react-router-dom";
import {Organizations} from "./components/organizations/Organizations";
import {OrganizationPortal} from "./components/organizations/OrganizationPortal";
import {OrganizationEditor} from "./components/organizations/OrganizationEditor";
import {Clients} from "./components/clients/Clients";
import {ClientEditor} from "./components/clients/ClientEditor";


function App() {
    return (
        <Fragment>
            <div className={"d-flex flex-column vh-100 vw-100 bg-dark"}>
                <BrowserRouter>
                    <NavBar/>
                    <Switch>
                        <Route path={"/organizations/add"} component={OrganizationEditor}/>
                        <Route path={"/organizations/:organizationId"} component={OrganizationPortal}/>
                        <Route path={"/organizations"} component={Organizations}/>
                        <Route path={"/clients/add"} component={ClientEditor}/>
                        <Route path={"/clients/:clientId"} component={ClientEditor}/>
                        <Route path={"/clients"} component={Clients}/>
                    </Switch>
                </BrowserRouter>
            </div>
        </Fragment>
    );
}

export default App;

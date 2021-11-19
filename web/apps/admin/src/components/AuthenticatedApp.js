"use strict";
exports.__esModule = true;
var NavBar_1 = require("./navbar/NavBar");
var react_router_dom_1 = require("react-router-dom");
var OrganizationEditor_1 = require("./organizations/OrganizationEditor");
var OrganizationPortal_1 = require("./organizations/OrganizationPortal");
var Organizations_1 = require("./organizations/Organizations");
var ClientEditor_1 = require("./clients/ClientEditor");
var Clients_1 = require("./clients/Clients");
var react_1 = require("react");
var AuthenticatedApp = function () {
    return (<div className={"d-flex flex-column vh-100 vw-100 bg-dark"}>
            <NavBar_1.NavBar />
            <react_router_dom_1.Switch>
                <react_router_dom_1.Route path={"/organizations/add"} component={OrganizationEditor_1.OrganizationEditor}/>
                <react_router_dom_1.Route path={"/organizations/:organizationId"} component={OrganizationPortal_1.OrganizationPortal}/>
                <react_router_dom_1.Route path={"/organizations"} component={Organizations_1.Organizations}/>
                <react_router_dom_1.Route path={"/clients/add"} component={ClientEditor_1.ClientEditor}/>
                <react_router_dom_1.Route path={"/clients/:clientId"} component={ClientEditor_1.ClientEditor}/>
                <react_router_dom_1.Route path={"/clients"} component={Clients_1.Clients}/>
            </react_router_dom_1.Switch>
        </div>);
};
exports["default"] = AuthenticatedApp;

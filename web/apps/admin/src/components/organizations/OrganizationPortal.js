"use strict";
exports.__esModule = true;
exports.OrganizationPortal = void 0;
var react_1 = require("react");
var react_router_dom_1 = require("react-router-dom");
var hooks_1 = require("../../hooks/hooks");
var OrganizationOverview_1 = require("./OrganizationOverview");
var OrganizationSideBar_1 = require("./OrganizationSideBar");
var IdentityProviders_1 = require("./identityproviders/IdentityProviders");
var IdentityProviderEditor_1 = require("./identityproviders/IdentityProviderEditor");
var OrganizationPortal = function (props) {
    var match = (0, react_router_dom_1.useRouteMatch)();
    var organization = (0, hooks_1.usePathOrganization)();
    if (!organization) {
        return <react_1.Fragment />;
    }
    return (<div className={"flex-grow-1 d-flex flex-column"}>
            <div className={"py-2 ps-4 bg-darkula text-white"}>
                <h5 className={"p-0 m-2"}>{organization.name}</h5>
            </div>
            <div className={"d-flex flex-row flex-grow-1 mt-4 px-4"}>
                <div className={""}>
                    <OrganizationSideBar_1.OrganizationSideBar organization={organization}/>
                </div>
                <div className={"flex-grow-1 ps-4 pe-2"}>
                    <react_router_dom_1.Switch>
                        {addIdentityProvidersRoute(match, organization)}
                        {identityProviderRoute(match, organization)}
                        {identityProvidersRoute(match, organization)}
                        {overviewRoute(match, organization)}
                    </react_router_dom_1.Switch>
                </div>
            </div>
        </div>);
};
exports.OrganizationPortal = OrganizationPortal;
function identityProvidersRoute(m, organization) {
    return <react_router_dom_1.Route path={m.path + "/identity-providers"} render={function () { return (<IdentityProviders_1.IdentityProviders organization={organization}/>); }}/>;
}
function addIdentityProvidersRoute(m, organization) {
    return <react_router_dom_1.Route path={m.path + "/identity-providers/add"} render={function () { return (<IdentityProviderEditor_1.IdentityProviderEditor organization={organization}/>); }}/>;
}
function identityProviderRoute(m, organization) {
    return <react_router_dom_1.Route path={m.path + "/identity-providers/:id"} render={function (p) { return (<IdentityProviderEditor_1.IdentityProviderEditor id={p.match.params["id"]} organization={organization}/>); }}/>;
}
function overviewRoute(m, organization) {
    return <react_router_dom_1.Route exact path={"" + m.path} render={function () { return <OrganizationOverview_1.OrganizationOverview organization={organization}/>; }}/>;
}

"use strict";
exports.__esModule = true;
exports.IdentityProviders = void 0;
var react_1 = require("react");
var hooks_1 = require("../../../hooks/hooks");
var react_router_dom_1 = require("react-router-dom");
var SectionTitle_1 = require("../../sectiontitle/SectionTitle");
var IdentityProviders = function (props) {
    var idps = (0, hooks_1.useIdentityProviders)(props.organization.id);
    var match = (0, react_router_dom_1.useRouteMatch)();
    return (<div>

            <SectionTitle_1.SectionTitle className={"text-light"} title={"Identity Providers"}>
                <react_router_dom_1.Link className={"btn btn-success btn-sm"} to={match.path + "/add"}>Add Identity Provider</react_router_dom_1.Link>
            </SectionTitle_1.SectionTitle>

            <div className={"list-group list-group-darkula"}>
                {idps.map(function (idp) { return (<react_router_dom_1.Link className={"list-group-item list-group-item-action"} to={match.url + "/" + idp.id}>
                        {idp.name} <span className={"badge bg-dark font-monospace"}>{idp.emailDomain}</span>
                    </react_router_dom_1.Link>); })}
                {idps.length === 0
            ? <div className={"disabled list-group-item"}>No Identity Provider</div>
            : <react_1.Fragment />}
            </div>
        </div>);
};
exports.IdentityProviders = IdentityProviders;

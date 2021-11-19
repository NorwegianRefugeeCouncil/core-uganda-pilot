"use strict";
exports.__esModule = true;
exports.OrganizationSideBar = void 0;
var react_router_dom_1 = require("react-router-dom");
var OrganizationSideBar = function (props) {
    var match = (0, react_router_dom_1.useRouteMatch)();
    return (<div className={"list-group list-group-darkula"} style={{ width: "15rem" }}>
            <react_router_dom_1.NavLink activeClassName={"active"} exact={true} className={"list-group-item list-group-item-action"} to={"" + match.url}>
                Overview
            </react_router_dom_1.NavLink>

            <react_router_dom_1.NavLink activeClassName={"active"} className={"list-group-item list-group-item-action"} to={match.url + "/identity-providers"}>
                Identity Providers
            </react_router_dom_1.NavLink>

        </div>);
};
exports.OrganizationSideBar = OrganizationSideBar;

"use strict";
exports.__esModule = true;
exports.Clients = void 0;
var react_1 = require("react");
var SectionTitle_1 = require("../sectiontitle/SectionTitle");
var hooks_1 = require("../../hooks/hooks");
var react_router_dom_1 = require("react-router-dom");
var Clients = function (props) {
    var apiClient = (0, hooks_1.useApiClient)();
    var _a = (0, react_1.useState)([]), clients = _a[0], setClients = _a[1];
    (0, react_1.useEffect)(function () {
        if (!apiClient) {
            return;
        }
        apiClient.listOAuth2Clients({}).then(function (resp) {
            if (resp.response) {
                setClients(resp.response.items);
            }
        });
    }, [apiClient]);
    return (<div className={"container mt-3"}>
            <div className={"row"}>
                <div className={"col"}>
                    <div className={"card card-darkula"}>
                        <div className={"card-body"}>
                            <SectionTitle_1.SectionTitle title={"OAuth2 Clients"}>
                                <react_router_dom_1.Link to={"/clients/add"} className={"btn btn-sm btn-success"}>Add OAuth2 Client</react_router_dom_1.Link>
                            </SectionTitle_1.SectionTitle>
                            <div className={"list-group list-group-darkula"}>
                                {clients.length === 0 && <div className={"list-group-item"}>No Clients</div>}
                                {clients.map(function (c) { return <react_router_dom_1.Link to={"/clients/" + c.id} className={"list-group-item"}>{c.clientName}</react_router_dom_1.Link>; })}
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>);
};
exports.Clients = Clients;

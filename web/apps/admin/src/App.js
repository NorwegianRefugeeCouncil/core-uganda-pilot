"use strict";
exports.__esModule = true;
var react_1 = require("react");
require("./App.scss");
require("bootstrap/dist/css/bootstrap.css");
require("bootstrap-icons/font/bootstrap-icons.css");
require("bootstrap/dist/js/bootstrap.bundle.min");
var react_router_dom_1 = require("react-router-dom");
var AuthenticatedApp_1 = require("./components/AuthenticatedApp");
var auth_client_1 = require("@core/auth-client");
var client_1 = require("../src/client/client");
function App() {
    var _a, _b;
    var client = new client_1["default"]();
    return (<react_router_dom_1.BrowserRouter>
            <react_router_dom_1.Switch>
                <auth_client_1.AuthWrapper clientId={((_a = process === null || process === void 0 ? void 0 : process.env) === null || _a === void 0 ? void 0 : _a.REACT_APP_CLIENT_ID) || ''} scopes={['openid', 'profile', 'offline_access']} issuer={((_b = process === null || process === void 0 ? void 0 : process.env) === null || _b === void 0 ? void 0 : _b.REACT_APP_ISSUER) || ''} redirectUriSuffix={'app'} axiosInstance={client.axiosInstance}>
                    <AuthenticatedApp_1["default"] />
                </auth_client_1.AuthWrapper>
            </react_router_dom_1.Switch>
        </react_router_dom_1.BrowserRouter>);
}
exports["default"] = App;

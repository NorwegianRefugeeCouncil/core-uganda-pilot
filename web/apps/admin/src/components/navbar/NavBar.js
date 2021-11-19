"use strict";
exports.__esModule = true;
exports.NavBar = void 0;
var react_router_dom_1 = require("react-router-dom");
var NavBar = function (props) {
    return (<nav className="navbar navbar-expand-sm navbar-light bg-light">
            <div className="container-fluid">
                <a className="navbar-brand" href="#">Core Admin</a>
                <button className="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarNav" aria-controls="navbarNav" aria-expanded="false" aria-label="Toggle navigation">
                    <span className="navbar-toggler-icon"/>
                </button>
                <div className="collapse navbar-collapse" id="navbarNav">
                    <ul className="navbar-nav">
                        <li className="nav-item">
                            <react_router_dom_1.NavLink className={"nav-link"} activeClassName={"active"} to={"/organizations"}>Organizations</react_router_dom_1.NavLink>
                        </li>
                        <li className="nav-item">
                            <react_router_dom_1.NavLink className={"nav-link"} activeClassName={"active"} to={"/clients"}>Clients</react_router_dom_1.NavLink>
                        </li>
                    </ul>
                </div>
            </div>
        </nav>);
};
exports.NavBar = NavBar;

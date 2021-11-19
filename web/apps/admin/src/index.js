"use strict";
exports.__esModule = true;
var react_1 = require("react");
var react_dom_1 = require("react-dom");
require("./index.css");
var App_1 = require("./App");
var reportWebVitals_1 = require("./utils/reportWebVitals");
react_dom_1["default"].render(<react_1["default"].StrictMode>
    <App_1["default"] />
  </react_1["default"].StrictMode>, document.getElementById('root'));
// If you want to start measuring performance in your hooks, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
(0, reportWebVitals_1["default"])();

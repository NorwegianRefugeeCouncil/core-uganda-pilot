"use strict";
exports.__esModule = true;
exports.axiosInstance = void 0;
var api_client_1 = require("@core/api-client");
var react_1 = require("react");
var host_1 = require("../constants/host");
var axios_1 = require("axios");
exports.axiosInstance = axios_1["default"].create();
function useApiClient() {
    return (0, react_1.useMemo)(function () {
        return new api_client_1["default"](host_1["default"], exports.axiosInstance);
    }, [1]);
}
exports["default"] = useApiClient;

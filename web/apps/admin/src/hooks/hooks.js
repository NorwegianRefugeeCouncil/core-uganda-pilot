"use strict";
exports.__esModule = true;
exports.useFormValidation = exports.useIdentityProviders = exports.usePathOrganization = exports.useApiClient = void 0;
var react_1 = require("react");
var client_1 = require("../client/client");
var react_router_dom_1 = require("react-router-dom");
var classnames_1 = require("classnames");
function useApiClient() {
    return (0, react_1.useMemo)(function () {
        return new client_1["default"]();
    }, []);
}
exports.useApiClient = useApiClient;
function usePathOrganization() {
    var apiClient = useApiClient();
    var organizationId = (0, react_router_dom_1.useParams)().organizationId;
    var _a = (0, react_1.useState)(), organization = _a[0], setOrganization = _a[1];
    (0, react_1.useEffect)(function () {
        if (!organizationId) {
            return;
        }
        apiClient.getOrganization({ id: organizationId }).then(function (resp) {
            if (resp.response) {
                setOrganization(resp.response);
            }
        });
    }, [apiClient, organizationId]);
    return organization;
}
exports.usePathOrganization = usePathOrganization;
function useIdentityProviders(organizationId) {
    var apiClient = useApiClient();
    var _a = (0, react_1.useState)([]), idps = _a[0], setIdps = _a[1];
    (0, react_1.useEffect)(function () {
        if (!organizationId) {
            return;
        }
        apiClient.listIdentityProviders({ organizationId: organizationId }).then(function (resp) {
            if (resp.response) {
                setIdps(resp.response ? resp.response.items : []);
            }
        });
    }, [apiClient, organizationId]);
    return idps;
}
exports.useIdentityProviders = useIdentityProviders;
function useFormValidation(isNew, form) {
    var _a = form.formState, dirtyFields = _a.dirtyFields, errors = _a.errors, isSubmitted = _a.isSubmitted, touchedFields = _a.touchedFields;
    console.log(isNew, touchedFields);
    var fieldClasses = (0, react_1.useCallback)(function (field) {
        var cls = (0, classnames_1["default"])("form-control form-control-darkula");
        var touched = touchedFields[field];
        var dirty = dirtyFields[field];
        var hasError = !!errors[field];
        console.log(errors);
        if (isSubmitted || (isNew && touched) || (isNew && dirty)) {
            if (hasError) {
                return (0, classnames_1["default"])(cls, "is-invalid");
            }
        }
        return cls;
    }, [dirtyFields, errors, isNew, isSubmitted, touchedFields]);
    var fieldErrors = (0, react_1.useCallback)(function (field) {
        var touched = touchedFields[field];
        var dirty = dirtyFields[field];
        var hasError = errors[field];
        var err = errors[field];
        return hasError && (isSubmitted || (isNew && touched) || (isNew && dirty)) &&
            <div className={"invalid-feedback"}>
                {(err === null || err === void 0 ? void 0 : err.type) === "required" ? <span>This field is required</span> : <react_1.Fragment />}
                {(err === null || err === void 0 ? void 0 : err.type) === "pattern" ? <span>Invalid value</span> : <react_1.Fragment />}
            </div>;
    }, [touchedFields, dirtyFields, errors, isNew, isSubmitted]);
    return { fieldErrors: fieldErrors, fieldClasses: fieldClasses };
}
exports.useFormValidation = useFormValidation;

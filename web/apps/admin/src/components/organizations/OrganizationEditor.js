"use strict";
exports.__esModule = true;
exports.OrganizationEditor = void 0;
var react_1 = require("react");
var SectionTitle_1 = require("../sectiontitle/SectionTitle");
var react_hook_form_1 = require("react-hook-form");
var hooks_1 = require("../../hooks/hooks");
var classnames_1 = require("classnames");
var react_router_dom_1 = require("react-router-dom");
var OrganizationEditor = function (props) {
    var apiClient = (0, hooks_1.useApiClient)();
    var params = (0, react_router_dom_1.useParams)();
    var _a = (0, react_1.useState)(), id = _a[0], setId = _a[1];
    var _b = (0, react_1.useState)(true), isNew = _b[0], setIsNew = _b[1];
    var form = (0, react_hook_form_1.useForm)({ mode: "onChange" });
    var register = form.register, handleSubmit = form.handleSubmit, setValue = form.setValue, isSubmitting = form.formState.isSubmitting;
    var _c = (0, hooks_1.useFormValidation)(!id, form), fieldClasses = _c.fieldClasses, fieldErrors = _c.fieldErrors;
    (0, react_1.useEffect)(function () {
        if (!apiClient) {
            return;
        }
        if (!setValue) {
            return;
        }
        if (!params.organizationId) {
            setId(undefined);
            setIsNew(true);
        }
        else if (params.organizationId) {
            setId(params.organizationId);
            setIsNew(false);
            apiClient
                .getIdentityProvider({ id: params.organizationId })
                .then(function (resp) {
                if (resp.response) {
                    setValue("name", resp.response.name);
                }
            });
        }
    }, [apiClient, setValue, params.organizationId]);
    var onSubmit = (0, react_1.useCallback)(function (args) {
        var obj = {
            id: id,
            name: args.name
        };
        if (id) {
            // return apiClient.updateIdentityProvider({
            //     object: obj
            // })
        }
        else {
            return apiClient.createOrganization({
                object: obj
            });
        }
    }, [apiClient, id]);
    return (<div className={"container mt-3"}>
            <div className={"row"}>
                <div className={"col"}>
                    <div className={"card card-darkula"}>
                        <div className={"card-body"}>
                            <SectionTitle_1.SectionTitle title={"Add Organization"}/>
                            <form autoComplete={"off"} onSubmit={handleSubmit(onSubmit)}>
                                <input type={"hidden"} autoComplete={"false"}/>
                                <div className={"form-group mb-2"}>
                                    <label className={"form-label"}>Name</label>
                                    <input type={"text"} className={(0, classnames_1["default"])("form-control form-control-darkula", fieldClasses("name"))} {...register("name", {
        required: true,
        pattern: /^[A-Za-z0-9_-]+( [A-Za-z0-9_-]+)*$/
    })}/>
                                    {fieldErrors("name")}
                                </div>
                                <div className={"py-2"}>
                                    <button disabled={isSubmitting} className={"btn btn-success"}>
                                        {isNew ? "Create Organization" : "Update Organization"}
                                    </button>
                                </div>
                            </form>
                        </div>
                    </div>
                </div>
            </div>
        </div>);
};
exports.OrganizationEditor = OrganizationEditor;

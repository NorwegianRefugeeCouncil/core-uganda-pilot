"use strict";
exports.__esModule = true;
exports.IdentityProviderEditor = void 0;
var react_1 = require("react");
var react_hook_form_1 = require("react-hook-form");
var classnames_1 = require("classnames");
var hooks_1 = require("../../../hooks/hooks");
var IdentityProviderEditor = function (props) {
    var id = props.id, organization = props.organization;
    var isNew = (0, react_1.useMemo)(function () { return !id; }, [id]);
    var apiClient = (0, hooks_1.useApiClient)();
    var form = (0, react_hook_form_1.useForm)({ mode: "onChange" });
    var register = form.register, handleSubmit = form.handleSubmit, setValue = form.setValue, isSubmitting = form.formState.isSubmitting;
    var _a = (0, hooks_1.useFormValidation)(isNew, form), fieldErrors = _a.fieldErrors, fieldClasses = _a.fieldClasses;
    (0, react_1.useEffect)(function () {
        if (id) {
            apiClient.getIdentityProvider({ id: id }).then(function (resp) {
                if (resp.response) {
                    setValue("name", resp.response.name);
                    setValue("clientId", resp.response.clientId);
                    setValue("organizationId", resp.response.organizationId);
                    setValue("issuer", resp.response.domain);
                    setValue("emailDomain", resp.response.emailDomain);
                    setValue("clientSecret", "");
                }
            });
        }
    }, [apiClient, id, setValue]);
    var onSubmit = (0, react_1.useCallback)(function (args) {
        var obj = {
            id: id,
            name: args.name,
            clientId: args.clientId,
            clientSecret: args.clientSecret,
            domain: args.issuer,
            organizationId: organization.id,
            emailDomain: args.emailDomain
        };
        if (id) {
            return apiClient.updateIdentityProvider({
                object: obj
            });
        }
        else {
            return apiClient.createIdentityProvider({
                object: obj
            });
        }
    }, [apiClient, id, organization.id]);
    return (<div className={(0, classnames_1["default"])("card bg-dark border-secondary")}>
            <div className={"card-body"}>
                <form className={"needs-validation"} noValidate onSubmit={handleSubmit(onSubmit)}>
                    <div className={(0, classnames_1["default"])("form-group mb-2")}>
                        <label className={"form-label text-light"}>Name</label>
                        <input {...register("name", {
        required: true,
        pattern: /^[a-zA-Z0-9\-_ ]+$/
    })} className={(0, classnames_1["default"])("form-control form-control-darkula", fieldClasses("name"))}/>
                        {fieldErrors("name")}
                    </div>
                    <div className={"form-group mb-2"}>
                        <label className={"form-label text-light"}>Issuer</label>
                        <input {...register("issuer", {
        required: true,
        pattern: /^https?:\/\/[a-zA-Z0-9.\-_]+(:[0-9]+)?$/
    })} className={(0, classnames_1["default"])("form-control form-control-darkula", fieldClasses("issuer"))}/>
                        {fieldErrors("issuer")}
                    </div>
                    <div className={"form-group mb-2"}>
                        <label className={"form-label text-light"}>Email Domain</label>
                        <input {...register("emailDomain", {
        required: true
    })} className={(0, classnames_1["default"])("form-control form-control-darkula", fieldClasses("emailDomain"))}/>
                        {fieldErrors("emailDomain")}
                    </div>
                    <div className={"form-group mb-2"}>
                        <label className={"form-label text-light"}>Client ID</label>
                        <input {...register("clientId", {
        required: true
    })} className={(0, classnames_1["default"])("form-control form-control-darkula", fieldClasses("clientId"))}/>
                        {fieldErrors("clientId")}
                    </div>
                    <div className={"form-group mb-2"}>
                        <label className={"form-label text-light"}>Client Secret</label>
                        <input type={"password"} {...register("clientSecret", {
        required: isNew
    })} className={(0, classnames_1["default"])("form-control form-control-darkula", fieldClasses("clientSecret"))} placeholder={isNew ? "" : "********"}/>
                        {fieldErrors("clientSecret")}
                    </div>
                    <button disabled={isSubmitting} className={"btn btn-success mt-2"}>
                        {props.id ? "Update Identity Provider" : "Create Identity Provider"}
                    </button>
                </form>
            </div>
        </div>);
};
exports.IdentityProviderEditor = IdentityProviderEditor;

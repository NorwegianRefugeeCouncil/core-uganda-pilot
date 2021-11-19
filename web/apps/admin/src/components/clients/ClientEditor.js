"use strict";
var __assign = (this && this.__assign) || function () {
    __assign = Object.assign || function(t) {
        for (var s, i = 1, n = arguments.length; i < n; i++) {
            s = arguments[i];
            for (var p in s) if (Object.prototype.hasOwnProperty.call(s, p))
                t[p] = s[p];
        }
        return t;
    };
    return __assign.apply(this, arguments);
};
exports.__esModule = true;
exports.ClientEditor = void 0;
var react_1 = require("react");
var SectionTitle_1 = require("../sectiontitle/SectionTitle");
var FormControl_1 = require("../formcontrol/FormControl");
var react_hook_form_1 = require("react-hook-form");
var hooks_1 = require("../../hooks/hooks");
var react_router_dom_1 = require("react-router-dom");
var ClientEditor = function (props) {
    var form = (0, react_hook_form_1.useForm)({ mode: "onTouched" });
    var register = form.register, handleSubmit = form.handleSubmit, isSubmitSuccessful = form.formState.isSubmitSuccessful;
    var _a = (0, hooks_1.useFormValidation)(true, form), fieldClasses = _a.fieldClasses, fieldErrors = _a.fieldErrors;
    var apiClient = (0, hooks_1.useApiClient)();
    var clientId = (0, react_router_dom_1.useParams)().clientId;
    var _b = (0, react_1.useState)(""), redirect = _b[0], setRedirect = _b[1];
    var _c = (0, react_1.useState)(""), clientSecret = _c[0], setClientSecret = _c[1];
    var setClient = (0, react_1.useCallback)(function (args) {
        form.setValue("clientName", args.clientName);
        form.setValue("uri", args.uri);
        form.setValue("scope", args.scope);
        form.setValue("tokenEndpointAuthMethod", args.tokenEndpointAuthMethod);
        form.setValue("responseTypes", args.responseTypes);
        form.setValue("grantTypes", args.grantTypes);
        form.setValue("redirectUris", args.redirectUris);
        form.setValue("allowedCorsOrigins", args.allowedCorsOrigins);
        setClientSecret(args.clientSecret);
    }, [form]);
    // gets the oauth2 client
    (0, react_1.useEffect)(function () {
        clientId && apiClient.getOAuth2Client({ id: clientId })
            .then(function (result) {
            result.response && setClient(result.response);
        });
    }, [form, clientId, setClient]);
    // updates or creates the oauth2 client on submit
    var onSubmit = (0, react_1.useCallback)(function (args) {
        clientId
            ? apiClient.updateOAuth2Client({ object: __assign(__assign({}, args), { id: clientId }) })
                .then(function (resp) { return resp.response && setClient(resp.response); })
            : apiClient.createOAuth2Client({ object: args })
                .then(function (resp) { return resp.response && setClient(resp.response); });
    }, [clientId, apiClient, setClient]);
    // deletes the oauth2 client && redirects to list
    var onDelete = (0, react_1.useCallback)(function () { return clientId && apiClient.deleteOAuth2Client({ id: clientId })
        .then(function () {
        setRedirect("/clients");
    }); }, [apiClient, clientId]);
    return (<react_1.Fragment>
            {redirect && <react_router_dom_1.Redirect to={redirect}/>}
            <div className={"container mt-3"}>
                <div className={"row"}>
                    <div className={"col"}>
                        <div className={"card card-darkula shadow"}>
                            <div className={"card-body"}>
                                <SectionTitle_1.SectionTitle title={"Create OAuth2 Client"}/>

                                <form noValidate={true} autoComplete={"off"} onSubmit={handleSubmit(onSubmit)}>
                                    <input type={"hidden"} autoComplete={"false"}/>

                                    {clientSecret && (<div className={"p-3 border border-warning shadow my-3"}>
                                            <FormControl_1.FormControl label={"Client Secret"} value={clientSecret} sensitive={true} allowCopy={true}/>
                                            <p className={"text-warning"}>Store the client secret securely! You will not
                                                be able to retrieve this value!</p>
                                        </div>)}

                                    <FormControl_1.FormControl {...register("clientName", { required: true })} label={"Name"} className={fieldClasses("clientName")}>
                                        {fieldErrors("clientName")}
                                    </FormControl_1.FormControl>

                                    <FormControl_1.FormControl {...register("uri", { required: true })} label={"URI"} className={fieldClasses("uri")}>
                                        {fieldErrors("uri")}
                                    </FormControl_1.FormControl>

                                    <FormControl_1.FormControl {...register("scope", { required: true })} label={"Scope"} className={fieldClasses("scope")}>
                                        {fieldErrors("scope")}
                                    </FormControl_1.FormControl>

                                    <FormControl_1.FormControl {...register("tokenEndpointAuthMethod", {
        required: true
    })} options={[
            { disabled: true, value: "", label: "Select token endpoint auth method" },
            { value: "client_secret_post", label: "Client Secret (post)" },
            { value: "client_secret_basic", label: "Client Secret (basic)" },
            { value: "private_key_jwt", label: "Private Key JWT" },
            { value: "none", label: "None" },
        ]} label={"Token Endpoint Auth Method"} className={fieldClasses("tokenEndpointAuthMethod")}>
                                        {fieldErrors("tokenEndpointAuthMethod")}
                                    </FormControl_1.FormControl>

                                    <FormControl_1.FormControl {...register("responseTypes", {
        required: true
    })} options={[
            { value: "id_token", label: "ID Token" },
            { value: "token", label: "Token" },
            { value: "code", label: "Code" },
        ]} multiple={true} label={"Response Types"} className={fieldClasses("responseTypes")}>
                                        {fieldErrors("responseTypes")}
                                    </FormControl_1.FormControl>

                                    <FormControl_1.FormControl {...register("grantTypes", {
        required: true
    })} options={[
            { value: "authorization_code", label: "Authorization Code" },
            { value: "refresh_token", label: "Refresh Token" },
            { value: "client_credentials", label: "Client Credentials" },
            { value: "implicit", label: "Implicit" },
        ]} multiple={true} label={"Grant Types"} className={fieldClasses("grantTypes")}>
                                        {fieldErrors("grantTypes")}
                                    </FormControl_1.FormControl>

                                    <FormControl_1.FormControl {...register("redirectUris", {
        required: true,
        setValueAs: function (value) {
            if (Array.isArray(value)) {
                return value;
            }
            return value
                .split(",")
                .map(function (v) { return v.trim(); })
                .filter(function (v) { return !!v; });
        }
    })} label={"Redirect URIs"} className={fieldClasses("redirectUris")}>
                                        {fieldErrors("redirectUris")}
                                    </FormControl_1.FormControl>

                                    <FormControl_1.FormControl {...register("allowedCorsOrigins", {
        required: true,
        setValueAs: function (value) {
            if (Array.isArray(value)) {
                return value;
            }
            return value
                .split(",")
                .map(function (v) { return v.trim(); })
                .filter(function (v) { return !!v; });
        }
    })} label={"Allowed CORS origins"} className={fieldClasses("allowedCorsOrigins")}>
                                        {fieldErrors("allowedCorsOrigins")}
                                    </FormControl_1.FormControl>


                                    <div className={"my-3"}>
                                        <button type={"submit"} className={"btn btn-success me-2"}>
                                            {clientId ? "Update OAuth2 Client" : "Create OAuth2 Client"}
                                        </button>
                                        {clientId && (<button onClick={onDelete} type={"button"} className={"btn btn-danger"}>
                                                Delete OAuth2 Client
                                            </button>)}
                                    </div>

                                    {isSubmitSuccessful && (<div className={"my-3 p-2 border border-success rounded text-success shadow"}>
                                            <i className={"bi bi-check"}/> Successfully {clientId ? "updated" : "created"} OAuth2
                                            client!
                                        </div>)}

                                </form>

                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </react_1.Fragment>);
};
exports.ClientEditor = ClientEditor;

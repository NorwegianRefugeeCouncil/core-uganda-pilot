import {FC, Fragment, useCallback, useEffect, useState} from "react";
import {SectionTitle} from "../sectiontitle/SectionTitle";
import {FormControl} from "../formcontrol/FormControl";
import {useForm} from "react-hook-form";
import {useApiClient, useFormValidation} from "../../app/hooks";
import {GrantType, OAuth2Client, ResponseType, TokenEndpointAuthMethod} from "../../client/client";
import {Redirect, useParams} from "react-router-dom";

export type FormData = {
    clientName: string
    uri: string
    grantTypes: GrantType[]
    responseTypes: ResponseType[]
    scope: string
    redirectUris: string[]
    allowedCorsOrigins: string[]
    tokenEndpointAuthMethod: TokenEndpointAuthMethod
    bla: string[]
}

export type RouteParams = {
    clientId?: string
}

export const ClientEditor: FC = props => {

    const form = useForm<FormData>({mode: "onTouched"})
    const {register, handleSubmit, formState: {isSubmitSuccessful}} = form
    const {fieldClasses, fieldErrors} = useFormValidation(true, form)
    const apiClient = useApiClient()
    const {clientId} = useParams<RouteParams>()
    const [redirect, setRedirect] = useState("")
    const [clientSecret, setClientSecret] = useState("")

    const setClient = useCallback((args: OAuth2Client) => {
        form.setValue("clientName", args.clientName)
        form.setValue("uri", args.uri)
        form.setValue("scope", args.scope)
        form.setValue("tokenEndpointAuthMethod", args.tokenEndpointAuthMethod)
        form.setValue("responseTypes", args.responseTypes)
        form.setValue("grantTypes", args.grantTypes)
        form.setValue("redirectUris", args.redirectUris)
        form.setValue("allowedCorsOrigins", args.allowedCorsOrigins)
        setClientSecret(args.clientSecret)
    }, [form])

    // gets the oauth2 client
    useEffect(() => {
        clientId && apiClient.getOAuth2Client({id: clientId})
            .then(result => result.response && setClient(result.response))
    }, [form, clientId, setClient])

    // updates or creates the oauth2 client on submit
    const onSubmit = useCallback((args: FormData) => {
        clientId
            ? apiClient.updateOAuth2Client({object: {...args, id: clientId}})
                .then(resp => resp.response && setClient(resp.response))
            : apiClient.createOAuth2Client({object: args})
                .then(resp => resp.response && setClient(resp.response))
    }, [clientId, apiClient, setClient])

    // deletes the oauth2 client && redirects to list
    const onDelete = useCallback(() => clientId && apiClient.deleteOAuth2Client({id: clientId})
            .then(() => {
                setRedirect("/clients")
            }),
        [apiClient, clientId])

    return (
        <Fragment>
            {redirect && <Redirect to={redirect}/>}
            <div className={"container mt-3"}>
                <div className={"row"}>
                    <div className={"col"}>
                        <div className={"card card-darkula shadow"}>
                            <div className={"card-body"}>
                                <SectionTitle title={"Create OAuth2 Client"}/>

                                <form noValidate={true} autoComplete={"off"} onSubmit={handleSubmit(onSubmit)}>
                                    <input type={"hidden"} autoComplete={"false"}/>

                                    {clientSecret && (
                                        <div className={"p-3 border border-warning shadow my-3"}>
                                            <FormControl
                                                label={"Client Secret"}
                                                value={clientSecret}
                                                sensitive={true}
                                                allowCopy={true}
                                            />
                                            <p className={"text-warning"}>Store the client secret securely! You will not
                                                be able to retrieve this value!</p>
                                        </div>
                                    )}

                                    <FormControl
                                        {...register("clientName", {required: true})}
                                        label={"Name"}
                                        className={fieldClasses("clientName")}>
                                        {fieldErrors("clientName")}
                                    </FormControl>

                                    <FormControl
                                        {...register("uri", {required: true})}
                                        label={"URI"}
                                        className={fieldClasses("uri")}>
                                        {fieldErrors("uri")}
                                    </FormControl>

                                    <FormControl
                                        {...register("scope", {required: true})}
                                        label={"Scope"}
                                        className={fieldClasses("scope")}>
                                        {fieldErrors("scope")}
                                    </FormControl>

                                    <FormControl
                                        {...register("tokenEndpointAuthMethod", {
                                            required: true,
                                        })}
                                        options={[
                                            {disabled: true, value: "", label: "Select token endpoint auth method"},
                                            {value: "client_secret_post", label: "Client Secret (post)"},
                                            {value: "client_secret_basic", label: "Client Secret (basic)"},
                                            {value: "private_key_jwt", label: "Private Key JWT"},
                                            {value: "none", label: "None"},
                                        ]}
                                        label={"Token Endpoint Auth Method"}
                                        className={fieldClasses("tokenEndpointAuthMethod")}>
                                        {fieldErrors("tokenEndpointAuthMethod")}
                                    </FormControl>

                                    <FormControl
                                        {...register("responseTypes", {
                                            required: true,
                                        })}
                                        options={[
                                            {value: "id_token", label: "ID Token"},
                                            {value: "token", label: "Token"},
                                            {value: "code", label: "Code"},
                                        ]}
                                        multiple={true}
                                        label={"Response Types"}
                                        className={fieldClasses("responseTypes")}>
                                        {fieldErrors("responseTypes")}
                                    </FormControl>

                                    <FormControl
                                        {...register("grantTypes", {
                                            required: true,
                                        })}
                                        options={[
                                            {value: "authorization_code", label: "Authorization Code"},
                                            {value: "refresh_token", label: "Refresh Token"},
                                            {value: "client_credentials", label: "Client Credentials"},
                                            {value: "implicit", label: "Implicit"},
                                        ]}
                                        multiple={true}
                                        label={"Grant Types"}
                                        className={fieldClasses("grantTypes")}>
                                        {fieldErrors("grantTypes")}
                                    </FormControl>

                                    <FormControl
                                        {...register("redirectUris", {
                                            required: true,
                                            setValueAs(value) {
                                                if (Array.isArray(value)) {
                                                    return value
                                                }
                                                return (value as string)
                                                    .split(",")
                                                    .map(v => v.trim())
                                                    .filter(v => !!v)
                                            }
                                        })}
                                        label={"Redirect URIs"}
                                        className={fieldClasses("redirectUris")}>
                                        {fieldErrors("redirectUris")}
                                    </FormControl>

                                    <FormControl
                                        {...register("allowedCorsOrigins", {
                                            required: true,
                                            setValueAs(value) {
                                                if (Array.isArray(value)) {
                                                    return value
                                                }
                                                return (value as string)
                                                    .split(",")
                                                    .map(v => v.trim())
                                                    .filter(v => !!v)
                                            }
                                        })}
                                        label={"Allowed CORS origins"}
                                        className={fieldClasses("allowedCorsOrigins")}>
                                        {fieldErrors("allowedCorsOrigins")}
                                    </FormControl>


                                    <div className={"my-3"}>
                                        <button type={"submit"} className={"btn btn-success me-2"}>
                                            {clientId ? "Update OAuth2 Client" : "Create OAuth2 Client"}
                                        </button>
                                        {clientId && (
                                            <button
                                                onClick={onDelete}
                                                type={"button"}
                                                className={"btn btn-danger"}>
                                                Delete OAuth2 Client
                                            </button>
                                        )}
                                    </div>

                                    {isSubmitSuccessful && (
                                        <div className={"my-3 p-2 border border-success rounded text-success shadow"}>
                                            <i className={"bi bi-check"}/> Successfully {clientId ? "updated" : "created"} OAuth2
                                            client!
                                        </div>
                                    )}

                                </form>

                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </Fragment>
    )

}

import {FC, useCallback} from "react";
import {SectionTitle} from "../sectiontitle/SectionTitle";
import {FormControl} from "../formcontrol/FormControl";
import {useForm} from "react-hook-form";
import {useApiClient, useFormValidation} from "../../app/hooks";
import {GrantType, ResponseType, TokenEndpointAuthMethod} from "../../client/client";

export type FormData = {
    name: string
    uri: string
    grantTypes: GrantType[]
    responseTypes: ResponseType[]
    scope: string
    redirectUris: string[]
    allowedCorsOrigins: string[]
    tokenEndpointAuthMethod: TokenEndpointAuthMethod
    bla: string[]
}

export const ClientEditor: FC = props => {

    const form = useForm<FormData>({mode: "onTouched"})
    const {register, handleSubmit} = form
    const {fieldClasses, fieldErrors} = useFormValidation(true, form)
    const apiClient = useApiClient()

    const onSubmit = useCallback((args: FormData) => {
        if (!apiClient) {
            return
        }
        apiClient.createOAuth2Client({object: args}).then(resp => {
            console.log(resp.response)
        })
    }, [apiClient])

    return (
        <div className={"container mt-3"}>
            <div className={"row"}>
                <div className={"col"}>
                    <div className={"card card-darkula"}>
                        <div className={"card-body"}>
                            <SectionTitle title={"Create OAuth2 Client"}/>

                            <form noValidate={true} autoComplete={"off"} onSubmit={handleSubmit(onSubmit)}>
                                <input type={"hidden"} autoComplete={"false"}/>

                                <FormControl
                                    {...register("name", {required: true})}
                                    label={"Name"}
                                    className={fieldClasses("name")}>
                                    {fieldErrors("name")}
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

                                <button className={"btn btn-success my-2"}>Create OAuth2 Client</button>

                            </form>

                        </div>
                    </div>
                </div>
            </div>
        </div>
    )

}

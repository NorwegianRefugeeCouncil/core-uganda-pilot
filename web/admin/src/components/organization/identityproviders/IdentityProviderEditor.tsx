import {FC, useCallback, useEffect} from "react";
import {useForm} from "react-hook-form";
import classNames from "classnames"
import {useApiClient} from "../../../app/hooks";
import {Organization} from "../../../client/client";

type Props = {
    id?: string
    organization: Organization
}

type FormData = {
    name: string,
    issuer: string,
    clientId: string,
    clientSecret: string,
    organizationId: string,
    emailDomain: string,
}

export const IdentityProviderEditor: FC<Props> = props => {

    const {
        id,
        organization
    } = props

    const apiClient = useApiClient()

    const {
        register,
        handleSubmit,
        setValue,
        formState: {errors, touchedFields, dirtyFields, isSubmitting, isSubmitted},
    } = useForm<FormData>({mode: "onChange"});

    useEffect(() => {
        if (id) {
            apiClient.getIdentityProvider({id}).then((resp) => {
                if (resp.response) {
                    setValue("name", resp.response.name)
                    setValue("clientId", resp.response.clientId)
                    setValue("organizationId", resp.response.organizationId)
                    setValue("issuer", resp.response.domain)
                    setValue("emailDomain", resp.response.emailDomain)
                    setValue("clientSecret", "")
                }
            })
        }
    }, [apiClient, id, setValue])

    const onSubmit = useCallback((args: FormData) => {
        let obj = {
            id: id,
            name: args.name,
            clientId: args.clientId,
            clientSecret: args.clientSecret,
            domain: args.issuer,
            organizationId: organization.id,
            emailDomain: args.emailDomain,
        };
        if (id) {
            return apiClient.updateIdentityProvider({
                object: obj
            })
        } else {
            return apiClient.createIdentityProvider({
                object: obj
            })
        }

    }, [apiClient, id, organization.id])

    const fieldClassNames = useCallback((field: keyof FormData) => {
        let cls = classNames("form-control form-control-darkula")
        if (isSubmitted || (!id && touchedFields[field]) || (!id && dirtyFields[field])) {
            if (errors[field]) {
                return classNames(cls, "is-invalid")
            }
        }
        return cls
    }, [dirtyFields, errors, id, isSubmitted, touchedFields])

    const fieldErrors = (field: keyof FormData) => {
        return errors[field] && (isSubmitted || (!id && touchedFields[field]) || (!id && dirtyFields[field])) &&
            <div className={"invalid-feedback"}>
                {errors[field]?.type === "required" && <span>This field is required</span>}
                {errors[field]?.type === "pattern" && <span>Invalid value</span>}
            </div>
    }

    return (
        <div className={classNames("card bg-dark border-secondary")}>
            <div className={"card-body"}>
                <form className={"needs-validation"} noValidate onSubmit={handleSubmit(onSubmit)}>
                    <div className={classNames("form-group mb-2")}>
                        <label className={"form-label text-light"}>Name</label>
                        <input
                            {...register("name", {
                                required: true,
                                pattern: /^[a-zA-Z0-9\-_ ]+$/,
                            })}
                            className={fieldClassNames("name")}
                        />
                        {fieldErrors("name")}
                    </div>
                    <div className={"form-group mb-2"}>
                        <label className={"form-label text-light"}>Issuer</label>
                        <input {...register("issuer", {
                            required: true,
                            pattern: /^https?:\/\/[a-zA-Z0-9.\-_]+(:[0-9]+)?$/,
                        })}
                               className={fieldClassNames("issuer")}
                        />
                        {fieldErrors("issuer")}
                    </div>
                    <div className={"form-group mb-2"}>
                        <label className={"form-label text-light"}>Email Domain</label>
                        <input {...register("emailDomain", {
                            required: true,
                        })}
                               className={fieldClassNames("emailDomain")}
                        />
                        {fieldErrors("emailDomain")}
                    </div>
                    <div className={"form-group mb-2"}>
                        <label className={"form-label text-light"}>Client ID</label>
                        <input {...register("clientId", {
                            required: true
                        })}
                               className={fieldClassNames("clientId")}
                        />
                        {fieldErrors("clientId")}
                    </div>
                    <div className={"form-group mb-2"}>
                        <label className={"form-label text-light"}>Client Secret</label>
                        <input
                            type={"password"}
                            {...register("clientSecret", {
                                required: !id
                            })}
                            className={fieldClassNames("clientSecret")}
                            placeholder={id ? "*****" : ""}
                        />
                        {fieldErrors("clientSecret")}
                    </div>
                    <button disabled={isSubmitting} className={"btn btn-success mt-2"}>
                        {props.id ? "Update Identity Provider" : "Create Identity Provider"}
                    </button>
                </form>
            </div>
        </div>
    )
}

import {FC, useCallback, useEffect, useState} from "react";
import {SectionTitle} from "../sectiontitle/SectionTitle";
import {useForm} from "react-hook-form";
import {useApiClient, useFormValidation} from "../../hooks/hooks";
import classNames from "classnames";
import {Organization} from "../../types/types";
import {useParams} from "react-router-dom";

export type  OrganizationEditorProps = {}

export type FormData = {
    name: string
}

export type RouteParams = {
    organizationId?: string
}

export const OrganizationEditor: FC<OrganizationEditorProps> = props => {

    const apiClient = useApiClient()
    const params = useParams<RouteParams>()
    const [id, setId] = useState<string>()
    const [isNew, setIsNew] = useState<boolean>(true)
    const form = useForm<FormData>({mode: "onChange"});
    const {
        register,
        handleSubmit,
        setValue,
        formState: {isSubmitting},
    } = form
    const {fieldClasses, fieldErrors} = useFormValidation(!id, form)

    useEffect(() => {
        if (!apiClient) {
            return
        }
        if (!setValue) {
            return
        }
        if (!params.organizationId) {
            setId(undefined)
            setIsNew(true)
        } else if (params.organizationId) {
            setId(params.organizationId)
            setIsNew(false)
            apiClient
                .getIdentityProvider({id: params.organizationId})
                .then((resp) => {
                    if (resp.response) {
                        setValue("name", resp.response.name)
                    }
                })
        }
    }, [apiClient, setValue, params.organizationId])

    const onSubmit = useCallback((args: FormData) => {
        let obj: Partial<Organization> = {
            id: id,
            name: args.name
        };
        if (id) {
            // return apiClient.updateIdentityProvider({
            //     object: obj
            // })
        } else {
            return apiClient.createOrganization({
                object: obj
            })
        }
    }, [apiClient, id])


    return (
        <div className={"container mt-3"}>
            <div className={"row"}>
                <div className={"col"}>
                    <div className={"card card-darkula"}>
                        <div className={"card-body"}>
                            <SectionTitle title={"Add Organization"}/>
                            <form autoComplete={"off"} onSubmit={handleSubmit(onSubmit)}>
                                <input type={"hidden"} autoComplete={"false"}/>
                                <div className={"form-group mb-2"}>
                                    <label className={"form-label"}>Name</label>
                                    <input
                                        type={"text"}
                                        className={classNames("form-control form-control-darkula", fieldClasses("name"))}
                                        {...register("name", {
                                            required: true,
                                            pattern: /^[A-Za-z0-9_-]+( [A-Za-z0-9_-]+)*$/
                                        })}
                                    />
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
        </div>
    )

}

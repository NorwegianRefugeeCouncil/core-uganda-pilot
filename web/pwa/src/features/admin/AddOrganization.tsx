import {FC, Fragment, useCallback, useState} from "react";
import {defaultClient} from "../../data/client";
import {Organization} from "../../types/types";
import {Redirect} from "react-router-dom"
import {useAppDispatch} from "../../app/hooks";
import {organizationActions} from "../../reducers/organizations";

export const AddOrganization: FC = props => {

    const dispatch = useAppDispatch()
    const [name, setName] = useState("")
    const [key, setKey] = useState("")
    const [createdOrganization,
        setCreatedOrganization] = useState<Organization | undefined>(undefined)
    const submit = useCallback(() => {
        defaultClient.createOrganization({object: {name, key}}).then((r) => {
            if (r.response) {
                dispatch(organizationActions.addOne(r.response))
                setCreatedOrganization(r.response)
            }
        })
    }, [dispatch, name, key])

    return (
        <div className={"col p-4"}>
            {createdOrganization ? <Redirect to={`/admin/organizations/${createdOrganization.id}`}/> : <Fragment/>}
            <div className={"form-group mb-2"}>
                <label className={"form-label"}>Key</label>
                <input className={"form-control"} value={key} onChange={ev => setKey(ev.target.value)}/>
            </div>
            <div className={"form-group mb-2"}>
                <label className={"form-label"}>Name</label>
                <input className={"form-control"} value={name} onChange={ev => setName(ev.target.value)}/>
            </div>
            <button className={"btn btn-primary"}
                    onClick={() => submit()}
                    disabled={key.length === 0 || name.length === 0}>
                Create Organization
            </button>
        </div>
    )

}

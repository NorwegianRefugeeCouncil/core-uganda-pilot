import React, {FC, Fragment, useCallback, useEffect, useState} from "react";
import {useAppDispatch, useAppSelector} from "../../app/hooks";
import {
    fetchIdentityProviders,
    identityProviderActions,
    identityProviderGlobalSelectors
} from "../../reducers/identityproviders";
import {IdentityProvider, Organization} from "../../types/types";
import {Redirect, useParams} from "react-router-dom"
import {defaultClient} from "../../data/client";

export type IdentityProviderProps = {
    id: string
    kind: string
    setKind: (kind: string) => void
    clientId: string
    setClientId: (clientId: string) => void
    clientSecret: string
    setClientSecret: (clientSecret: string) => void
    domain: string
    setDomain: (domain: string) => void
    canSave: boolean
    onSave: () => void
    success: boolean
}

const IdentityProviderEditor: FC<IdentityProviderProps> = (props) => {
    const {
        id,
        kind,
        setKind,
        clientId,
        setClientId,
        clientSecret,
        setClientSecret,
        domain,
        setDomain,
        canSave,
        onSave,
        success
    } = props
    return (
        <div className={"card text-dark"}>
            <div className={"card-body"}>
                <h3>Identity Provider</h3>
                <div className={"form-group mb-2"}>
                    <label className={"form-label"}>Kind</label>
                    <select defaultValue={""} value={kind} onChange={(e) => setKind(e.target.value)} className={"form-select"}>
                        <option value={""} disabled={true}/>
                        <option value={"oidc"}>OIDC</option>
                    </select>
                </div>
                <div className={"form-group mb-2"}>
                    <label className={"form-label"}>URL</label>
                    <input value={domain}
                           onChange={(e) => setDomain(e.target.value)}
                           className={"form-control font-monospace"}/>
                </div>
                <div className={"form-group mb-2"}>
                    <label className={"form-label"}>Client ID</label>
                    <input
                        value={clientId}
                        onChange={e => setClientId(e.target.value)}
                        className={"form-control font-monospace"}/>
                </div>
                <div className={"form-group mb-2"}>
                    <label className={"form-label"}>Client Secret</label>
                    <input value={clientSecret} onChange={e => setClientSecret(e.target.value)}
                           placeholder={id ? "************" : ""} type={"password"} className={"form-control"}/>
                </div>
                <button
                    disabled={!canSave}
                    onClick={() => onSave()}
                    className={"btn btn-primary my-2"}>
                    {id ? "Save Identity Provider" : "Add Identity Provider"}
                </button>
                {success
                    ? <p className={"text-success fw-bold"}>Successfully saved!</p>
                    : <Fragment/>}
            </div>
        </div>
    )
}

export type IdentityProviderEditorContainerProps = {
    organization: Organization
}

export const IdentityProviderEditorContainer: FC<IdentityProviderEditorContainerProps> = props => {

    const [pending, setPending] = useState(false)
    const [error, setError] = useState<any>(undefined)
    const [success, setSuccess] = useState(false)
    const [created, setCreated] = useState<IdentityProvider | undefined>(undefined)
    const [id, setId] = useState("")
    const [domain, setDomain] = useState("")
    const [kind, setKind] = useState("")
    const [clientId, setClientId] = useState("")
    const [clientSecret, setClientSecret] = useState("")
    const [canSave, setCanSave] = useState(false)

    const {identityProviderId} = useParams<{ identityProviderId: string | undefined }>()
    const organization = props.organization

    const dispatch = useAppDispatch()
    useEffect(() => {
        dispatch(fetchIdentityProviders({organizationId: organization.id}))
    }, [dispatch, organization.id])

    const idp = useAppSelector(state => {
        if (identityProviderId) {
            return identityProviderGlobalSelectors.selectById(state, identityProviderId)
        }
    })

    useEffect(() => {
        if (!idp) {
            setId("")
            setDomain("")
            setKind("")
            setClientId("")
            setClientSecret("")
        } else if (idp.id !== id) {
            setId(idp.id)
            setDomain(idp.domain)
            setKind(idp.kind)
            setClientId(idp.clientId)
            setClientSecret("")
        }
    }, [id, idp])

    useEffect(() => {
        if (!kind || !domain || !clientId) {
            setCanSave(false)
            return
        }
        if (!id && !clientSecret) {
            setCanSave(false)
            return
        }
        if (pending) {
            setCanSave(false)
            return
        }
        if (pending) {
            setCanSave(false)
            return
        }
        if (success) {
            setCanSave(false)
            return
        }
        setCanSave(true)
    }, [success, kind, domain, clientId, clientSecret, id, pending, created])

    const onSave = useCallback(async () => {
        const obj = {kind, domain, clientSecret, clientId, id, organizationId: organization.id}
        let promise: Promise<IdentityProvider | undefined>

        if (!id) {
            promise = defaultClient.createIdentityProvider({object: obj}).then(r => {
                if (r.error) {
                    throw r.error
                }
                return r.response
            })
        } else {
            promise = defaultClient.updateIdentityProvider({object: obj}).then(r => {
                if (r.error) {
                    throw r.error
                }
                return r.response
            })
        }

        try {
            const result = await promise
            setError(undefined)
            setCreated(result)
            if (result){
                setId(result.id)
                setKind(result.kind)
                setCreated(result)
                setDomain(result.domain)
                setClientSecret(result.clientSecret)
                setClientId(result.clientId)
                dispatch(identityProviderActions.upsertOne(result))
                setSuccess(true)
            }
        } catch (err) {
            setError(err)
        } finally {
            setPending(false)
        }

    }, [dispatch, kind, domain, clientSecret, clientId, id, organization.id])

    const renderPending = useCallback(() => {
        if (pending) {
            return <div className="spinner-border text-primary" role="status">
                <span className="sr-only">Loading...</span>
            </div>
        }
    }, [pending])

    const renderError = useCallback(() => {
        if (error) {
            return <code className={"mt-2"}>
                {JSON.stringify(error)}
            </code>
        }
    }, [error])


    useEffect(() => {
        if (idp){
            if (idp.kind !== kind){
                setSuccess(false)
                return
            }
            if (clientSecret){
                setSuccess(false)
                return
            }
            if (idp.clientId !== clientId){
                setSuccess(false)
                return
            }
            if (idp.domain !== domain){
                setSuccess(false)
                return
            }
            setSuccess(false)
        }
    }, [kind, clientSecret, clientId, domain, idp])

    if (id && id !== identityProviderId) {
        return <Redirect to={`/admin/organizations/${organization.id}/identityproviders/${id}`}/>
    }

    if (!organization) {
        return <Fragment/>
    }

    return (
        <div>
            {renderPending()}
            <IdentityProviderEditor
                id={id}
                onSave={onSave}
                setKind={setKind}
                clientId={clientId}
                kind={kind}
                domain={domain}
                canSave={canSave}
                clientSecret={clientSecret}
                setClientId={setClientId}
                setClientSecret={setClientSecret}
                setDomain={setDomain}
                success={success}
            />
            {renderError()}
        </div>
    )

}

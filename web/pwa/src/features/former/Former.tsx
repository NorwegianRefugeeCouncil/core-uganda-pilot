import React, {FC, Fragment, useCallback, useEffect, useState} from "react";
import {useAppDispatch, useAppSelector} from "../../app/hooks";
import {fetchDatabases} from "../../reducers/database";
import {fetchFolders} from "../../reducers/folder";
import {fetchForms} from "../../reducers/form";
import {useLocation} from "react-router-dom"
import {formerActions, formerGlobalSelectors, FormField, postForm} from "./former.slice";
import {FieldKind} from "../../types/types";
import {FormerField} from "./Field";
import {FormName} from "./FormName";
import {FieldTypePicker} from "./FieldTypePicker";


type FormerProps = {
    formName: string,
    setFormName: (formName: string) => void,
    fields: FormField[]
    fieldOptions?: string[]
    setFieldOption: (fieldId: string, i: number, value: string) => void
    addOption: (fieldId: string) => void
    removeOption: (fieldId: string, index: number) => void
    selectedFieldId: string | undefined
    setSelectedField: (fieldId: string | undefined) => void
    addField: (kind: FieldKind) => void
    setFieldRequired: (fieldId: string, required: boolean) => void
    setFieldIsKey: (fieldId: string, isKey: boolean) => void
    setFieldName: (fieldId: string, name: string) => void
    setFieldDescription: (fieldId: string, description: string) => void
    setFieldReferencedDatabaseId: (fieldId: string, databaseId: string) => void
    setFieldReferencedFormId: (fieldId: string, formId: string) => void
    openSubForm: (fieldId: string) => void
    saveField: (fieldId: string) => void
    saveForm: () => void
    parentFormName: string | undefined
    cancelField: (fieldId: string) => void
}

function mapField(f: FormField, props: FormerProps) {
    const {
        selectedFieldId,
        setSelectedField,
        setFieldName,
        setFieldOption,
        addOption,
        removeOption,
        setFieldDescription,
        setFieldIsKey,
        setFieldRequired,
        openSubForm,
        saveField,
        cancelField,
        setFieldReferencedDatabaseId,
        setFieldReferencedFormId
    } = props

    return <FormerField
        key={f.id}
        isSelected={f.id === selectedFieldId}
        selectField={() => setSelectedField(f.id)}
        fieldType={f.type}
        fieldOptions={f.options}
        setFieldOption={(i: number, value: string) => setFieldOption(f.id, i, value)}
        addOption={() => addOption(f.id)}
        removeOption={(i: number) => removeOption(f.id, i)}
        fieldName={f.name}
        setFieldName={(name) => setFieldName(f.id, name)}
        fieldRequired={f.required}
        setFieldRequired={(req) => setFieldRequired(f.id, req)}
        fieldIsKey={f.key}
        setFieldIsKey={(isKey) => setFieldIsKey(f.id, isKey)}
        fieldDescription={f.description}
        setFieldDescription={(d) => setFieldDescription(f.id, d)}
        openSubForm={() => openSubForm(f.id)}
        cancel={() => cancelField(f.id)}
        saveField={() => saveField(f.id)}
        referencedDatabaseId={f.referencedDatabaseId}
        referencedFormId={f.referencedFormId}
        setReferencedDatabaseId={(d) => setFieldReferencedDatabaseId(f.id, d)}
        setReferencedFormId={(d) => setFieldReferencedFormId(f.id, d)}
    />
}


export const Former: FC<FormerProps> = props => {

    const {
        formName,
        setFormName,
        fields,
        selectedFieldId,
        addField,
        saveForm,
        parentFormName
    } = props

    const [isAddingField, setIsAddingField] = useState(false)

    const selectedField = selectedFieldId ? fields.find(f => f.id === selectedFieldId) : undefined

    function formHeader() {
        return <FormName formName={formName} setFormName={setFormName}/>
    }

    function addFieldButton() {
        return <div>
            <button className={"btn btn-primary my-2 mb-3 w-100"}
                    onClick={() => setIsAddingField(true)}>
                Add field
            </button>
        </div>
    }


    const fieldSections = useCallback(() => {
        if (isAddingField) {
            return <FieldTypePicker
                onCancel={() => {
                    setIsAddingField(false)
                }}
                onSubmit={(fieldKind) => {
                    addField(fieldKind)
                    setIsAddingField(false)
                }}/>
        } else {
            return <div>
                {addFieldButton()}
                {fields.map((f) => mapField(f, props))}
            </div>
        }
    }, [addField, fields, isAddingField, props])


    if (selectedField) {
        return <div className={"flex-grow-1 overflow-scroll"}>
            <div className={"container-fluid mt-4"}>
                <div className={"row"}>
                    <div className={"col-12 col-md-8 offset-md-1"}>
                        <h3>Add Form</h3>
                        <h6>
                            {parentFormName
                                ? <div className={"mb-2"}>Editing child form of {parentFormName}</div>
                                : <Fragment/>}
                        </h6>
                    </div>
                </div>
                <div className={"row mt-3"}>
                    <div className={"col-10 col-md-8 offset-md-1"}>
                        {formHeader()}
                        {mapField(selectedField, props)}
                    </div>
                    <div className={"col-2"}>
                        <button className={"btn btn-primary w-100"} onClick={() => saveForm()}>Save</button>
                    </div>
                </div>
            </div>
        </div>
    }

    return (
        <div className={"h-100 w-100 flex-grow-1 overflow-scroll"}>
            <div className={"container mt-4"}>
                <div className={"row"}>
                    <div className={"col-8 offset-2"}>
                        <h3>Add Form</h3>
                        <h6>
                            {parentFormName
                                ? <div className={"mb-2 p-2 border-secondary"}>Editing child form
                                    of {parentFormName}</div>
                                : <Fragment/>}
                        </h6>
                    </div>
                </div>
                <div className={"row mt-3"}>
                    <div className={"col-6 offset-2"}>
                        {formHeader()}
                        {fieldSections()}
                    </div>
                    <div className={"col-2"}>
                        <button className={"btn btn-primary w-100"} onClick={() => saveForm()}>Save</button>
                    </div>
                </div>
            </div>
        </div>)

}

export const FormerContainer: FC = props => {

    const dispatch = useAppDispatch()

    // load data
    useEffect(() => {
        dispatch(formerActions.reset())
        dispatch(fetchDatabases())
        dispatch(fetchFolders())
        dispatch(fetchForms())
    }, [dispatch])

    const location = useLocation()

    const form = useAppSelector(formerGlobalSelectors.selectCurrentForm)
    const parentForm = useAppSelector(formerGlobalSelectors.selectCurrentFormParent)
    const folder = useAppSelector(formerGlobalSelectors.selectFolder)
    const database = useAppSelector(formerGlobalSelectors.selectDatabase)
    const selectedField = useAppSelector(formerGlobalSelectors.selectCurrentField)

    const formDefinition = useAppSelector(formerGlobalSelectors.selectFormDefinition(database?.id, folder?.id))

    useEffect(() => {
        const search = new URLSearchParams(location.search)
        let databaseId = search.get("databaseId");
        if (databaseId) {
            dispatch(formerActions.setDatabase({databaseId: databaseId}))
        }
        const folderId = search.get("folderId")
        if (folderId) {
            dispatch(formerActions.setFolder({folderId: folderId}))
        }
    }, [dispatch, location])

    const setFormName = useCallback((formName: string) => {
        if (form) {
            dispatch(formerActions.setFormName({formId: form.formId, formName: formName}))
        }
    }, [dispatch, form])

    const setSelectedField = useCallback((fieldId: string | undefined) => {
        dispatch(formerActions.selectField({fieldId}))
    }, [dispatch])

    const addField = useCallback((kind: FieldKind) => {
        if (form) {
            dispatch(formerActions.addField({formId: form.formId, kind: kind}))
        }
    }, [dispatch, form])

    const setFieldOption = useCallback((fieldId: string, i: number, value: string) => {
        dispatch(formerActions.setFieldOption({fieldId, i, value}))
    }, [dispatch])

    const addOption = useCallback((fieldId: string) => {
        dispatch(formerActions.addOption({fieldId}))
    }, [dispatch])

    const removeOption = useCallback((fieldId:string, i: number) => {
        dispatch(formerActions.removeOption({fieldId, i}))
    }, [dispatch])

    const cancelField = useCallback((fieldId: string) => {
        dispatch(formerActions.cancelFieldChanges({fieldId}))
    }, [dispatch])

    const setFieldRequired = useCallback((fieldId: string, required: boolean) => {
        dispatch(formerActions.setFieldRequired({fieldId, required}))
    }, [dispatch])

    const setFieldIsKey = useCallback((fieldId: string, isKey: boolean) => {
        dispatch(formerActions.setFieldIsKey({fieldId, isKey}))
    }, [dispatch])

    const setFieldName = useCallback((fieldId: string, name: string) => {
        dispatch(formerActions.setFieldName({fieldId, name}))
    }, [dispatch])

    const setFieldDescription = useCallback((fieldId: string, description: string) => {
        dispatch(formerActions.setFieldDescription({fieldId, description}))
    }, [dispatch])

    const setFieldReferencedDatabaseId = useCallback((fieldId: string, databaseId: string) => {
        dispatch(formerActions.setFieldReferencedDatabaseId({fieldId, databaseId}))
    }, [dispatch])

    const setFieldRefernecedFormId = useCallback((fieldId: string, formId: string) => {
        dispatch(formerActions.setFieldReferencedFormId({fieldId, formId}))
    }, [dispatch])

    const openSubForm = useCallback((fieldId: string) => {
        dispatch(formerActions.openSubForm({fieldId: fieldId}))
    }, [dispatch])

    const saveField = useCallback((fieldId: string) => {
        dispatch(formerActions.selectField({fieldId: undefined}))
    }, [dispatch])

    const saveForm = useCallback(() => {
        if (parentForm) {
            dispatch(formerActions.saveForm())
        } else {
            if (formDefinition) {
                dispatch(postForm(formDefinition))
            }
        }

    }, [dispatch, formDefinition, parentForm])

    if (!form) {
        return <Fragment/>
    }


    return <Former
        formName={form.name}
        setFormName={setFormName}
        fields={form.fields}
        selectedFieldId={selectedField?.id}
        setSelectedField={setSelectedField}
        addField={addField}
        setFieldOption={setFieldOption}
        addOption={addOption}
        removeOption={removeOption}
        setFieldRequired={setFieldRequired}
        setFieldIsKey={setFieldIsKey}
        setFieldName={setFieldName}
        setFieldDescription={setFieldDescription}
        openSubForm={openSubForm}
        saveField={saveField}
        saveForm={saveForm}
        parentFormName={parentForm?.name}
        cancelField={(fieldId: string) => cancelField(fieldId)}
        setFieldReferencedDatabaseId={setFieldReferencedDatabaseId}
        setFieldReferencedFormId={setFieldRefernecedFormId}
    />

}

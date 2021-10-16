import React, {Fragment, useCallback, useEffect, useState} from "react";
import {Link, Redirect} from "react-router-dom";
import {FieldDefinition, FieldType, FieldTypeSubForm, FormDefinition} from "./client";
import {
    createForm,
    createFormAddField,
    createFormReset, createFormSaveField,
    createFormSetName, createFormSetOpenSubForm, createFormSetSaveSubForm,
    createFormSetSelectedField,
    createFormSetSelectedFieldCode,
    createFormSetSelectedFieldDescription,
    createFormSetSelectedFieldIsKey,
    createFormSetSelectedFieldName,
    createFormSetSelectedFieldRequired,
    createFormSetSelectedFieldType,
    fetchDatabases,
    fetchFolders,
    fieldPropsInitialState,
    FieldPropsState,
    FormPropsState,
    store
} from "./store";
import {FieldEditor, GetFieldText} from "./FieldEditor";
import {useSelectedDatabase, useSelectedFolders} from "./utils";
import {Header} from "./Header";


export function FormEditor() {

    const database = useSelectedDatabase()
    const selectedFolders = useSelectedFolders()

    const formName = useEditFormName()
    const selectedField = useEditFormSelectedField()
    const fields = useEditFormFields()
    const allFields = useEditFormAllFields()
    const forms = useEditFormForms()
    const parentForm = useParentForm()

    const [createFormId, setCreateFormId] = useState<string | null>("")

    useEffect(() => {
        store.dispatch(fetchDatabases())
        store.dispatch(fetchFolders({options: {}}))
        const sub = store.state$.subscribe(s => {
            setCreateFormId(s.form.createdFormId)
        })
        return () => {
            sub.unsubscribe()
        }
    }, [])

    const redirectIfDone = useCallback(() => {
        if (createFormId) {
            setTimeout(() => {
                store.dispatch(createFormReset())
            })
            const redirectTo = selectedFolders.length > 0
                ? `/databases/${database?.id}/forms/${createFormId}?folderId=${selectedFolders[0].id}`
                : `/databases/${database?.id}/forms/${createFormId}`
            return <Redirect to={redirectTo}/>
        }
        return <Fragment/>
    }, [database, createFormId, selectedFolders])

    const parseFields = (formIdx : number, flds: FieldPropsState[]): FieldDefinition[] => {
        const result: FieldDefinition[] = []
        for (let i = 0; i < flds.length; i++) {
            const field = flds[i]
            if (field.formIdx !== formIdx) {
                continue
            }
            let fieldType = new FieldType()
            if (field.type === "text") {
                fieldType.text = {}
            } else if (field.type === "subform") {
                const subForm = forms[field.formIdx]
                fieldType.subForm = new FieldTypeSubForm()
                fieldType.subForm.code = subForm.code
                fieldType.subForm.name = subForm.name
                fieldType.subForm.fields = parseFields(field.subFormIdx, flds)
            }
            result.push({
                name: field.name,
                description: field.description,
                id: "",
                code: field.code,
                required: field.isRequired,
                fieldType: fieldType,
            })
        }
        return result
    }

    const submitForm = () => {
        if (parentForm) {
            store.dispatch(createFormSetSaveSubForm())
        } else {
            const formDef = new FormDefinition()
            formDef.name = forms[0].name
            formDef.databaseId = database?.id as string
            formDef.folderId = selectedFolders.length > 0
                ? selectedFolders[0].id
                : ""
            formDef.fields = parseFields(0, allFields)
            store.dispatch(createForm({form: formDef}))
        }
    }

    const setSelectedField = (index: number) => {
        store.dispatch(createFormSetSelectedField({index}))
    }

    const addField = () => {
        store.dispatch(createFormAddField({field: fieldPropsInitialState, select: true}))
    }

    const handleSaveFieldDefinition = () => {
        if (!selectedField) {
            return
        }
        store.dispatch(createFormSaveField({index: selectedField?.index}))
        setSelectedField(-1)
    }

    if (!database) {
        return <Fragment/>
    }

    return <Fragment>
        {createFormId}
        {redirectIfDone()}
        <Header
            database={database}
            title={"Create Database"}
            folders={selectedFolders}
            additionalBreadcrumb={"New Form"}>
        </Header>
        <main className={"bg-dark text-light vh-100 py-5"}>
            <div className={"container"}>
                <div className={"row"}>
                    <div className={"col-10"}>
                        <h5>Add Form</h5>

                        {!parentForm
                            ? <Fragment/>
                            : <div>
                                Sub Form of {parentForm.name}
                            </div>}

                        <div className={"form-group"}>
                            <label className={"form-label"} htmlFor={"name"}>Form Name:</label>
                            <input
                                onChange={e => store.dispatch(createFormSetName({name: e.target.value}))}
                                value={formName ? formName : ""}
                                className={"form-control"}
                                name={"name"}/>
                        </div>


                        <div className={"mt-3"}>
                            {!selectedField
                                ? <div className={"text-center mb-4"}>
                                    <button
                                        onClick={() => addField()}
                                        className={"btn btn-primary mt-3"}>
                                        Add field
                                    </button>
                                </div>
                                : <Fragment/>}
                        </div>

                        {fields.map(f => {

                            if (f.index === selectedField?.index) {
                                return <div className={"mb-2"}>
                                    <FieldEditor
                                        index={selectedField.index}
                                        formIdx={selectedField.formIdx}
                                        subFormIdx={selectedField.subFormIdx}
                                        code={selectedField.code}
                                        setCode={(c) => store.dispatch(createFormSetSelectedFieldCode({code: c}))}
                                        isKey={selectedField.isKey}
                                        setIsKey={(k) => store.dispatch(createFormSetSelectedFieldIsKey({isKey: k}))}
                                        description={selectedField.description}
                                        setDescription={(d) => store.dispatch(createFormSetSelectedFieldDescription({description: d}))}
                                        name={selectedField.name}
                                        setName={(n) => store.dispatch(createFormSetSelectedFieldName({name: n}))}
                                        type={selectedField.type}
                                        setKind={kind => store.dispatch(createFormSetSelectedFieldType({kind}))}
                                        isRequired={selectedField.isRequired}
                                        setIsRequired={required => store.dispatch(createFormSetSelectedFieldRequired({required}))}
                                        onSave={() => handleSaveFieldDefinition()}
                                        onOpenSubForm={() => store.dispatch(createFormSetOpenSubForm())}
                                    />
                                </div>
                            }

                            return <Fragment key={f.name}>
                                <div
                                    onClick={() => store.dispatch(createFormSetSelectedField({index: f.index}))}
                                    style={{cursor: "pointer"}}
                                    className={"p-3 d-flex flex-row border border-secondary mb-2"}>
                                    <div>
                                        <div className={"text-muted"}>
                                            {GetFieldText({fieldDef: f})}
                                        </div>
                                        <div>
                                            {f.name}
                                        </div>
                                    </div>
                                </div>
                            </Fragment>
                        })}
                    </div>

                    <div className={"col-2"}>
                        <div className={"d-flex flex-column"}>
                            <button
                                className={"btn btn-sm btn-primary"}
                                onClick={() => submitForm()}
                                disabled={fields.length === 0 || !formName}>
                                Save
                            </button>
                            <button
                                className={"btn btn-sm btn-outline-primary mt-2"}>
                                Cancel
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        </main>
    </Fragment>

}


function useEditFormName() {
    const [formName, setFormName] = useState<string>("")
    useEffect(() => {
        const sub = store.createFormName$.subscribe(name => {
            setFormName(name)
        })
        return () => {
            sub.unsubscribe()
        }
    }, [])
    return formName
}


function useEditFormForms() {
    const [forms, setForms] = useState<FormPropsState[]>([])
    useEffect(() => {
        const sub = store.createFormForms$.subscribe(f => {
            setForms(f)
        })
        return () => {
            sub.unsubscribe()
        }
    }, [])
    return forms
}

function useIsSubForm() {
    const [isSubForm, setIsSubForm] = useState<boolean>(false)
    useEffect(() => {
        const sub = store.createFormSelectedFormIsSubForm$.subscribe(f => {
            setIsSubForm(f)
        })
        return () => {
            sub.unsubscribe()
        }
    }, [])
    return isSubForm
}

function useParentForm() {
    const [parentForm, setParentForm] = useState<FormPropsState | undefined>(undefined)
    useEffect(() => {
        const sub = store.createFormParentForm$.subscribe(f => {
            setParentForm(f)
        })
        return () => {
            sub.unsubscribe()
        }
    }, [])
    return parentForm
}

function useEditFormFields() {
    const [formFields, setFormFields] = useState<FieldPropsState[]>([])
    useEffect(() => {
        const sub = store.createFormFields$.subscribe(fields => {
            setFormFields(fields)
        })
        return () => {
            sub.unsubscribe()
        }
    }, [])
    return formFields
}

function useEditFormAllFields() {
    const [formFields, setFormFields] = useState<FieldPropsState[]>([])
    useEffect(() => {
        const sub = store.createFormAllFields$.subscribe(fields => {
            setFormFields(fields)
        })
        return () => {
            sub.unsubscribe()
        }
    }, [])
    return formFields
}

function useEditFormSelectedField() {
    const [selectedField, setSelectedField] = useState<FieldPropsState | undefined>(undefined)
    useEffect(() => {
        const sub = store.createFormSelectedField$.subscribe(field => {
            setSelectedField(field)
        })
        return () => {
            sub.unsubscribe()
        }
    }, [])
    return selectedField
}


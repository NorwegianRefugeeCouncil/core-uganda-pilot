import {
    useSelectedDatabase,
    useSelectedForm,
    usePostRecordSuccess,
    useSelectedRecord,
    useSelectedFormIntf
} from "./utils";
import {FC, Fragment, useEffect, useRef, useState} from "react";
import {Header} from "./Header";
import {Database, FieldDefinition, FieldTypeSubForm, FormDefinition} from "./client";
import {
    createRecord,
    createRecordCancel,
    fetchDatabases,
    fetchForms, FormIntf,
    initNewRecord, recordOpenSubForm,
    setRecordValue, setSelectedForm,
    store
} from "./store";
import {Redirect} from "react-router-dom";

export type RecordEditorContainerProps = {}

export function RecordEditorContainer() {

    const selectedDatabase = useSelectedDatabase()
    const selectedForm = useSelectedForm()
    const selectedRecord = useSelectedRecord()
    const postRecordSuccess = usePostRecordSuccess()
    const selectedFormIntf = useSelectedFormIntf()

    const values = useRef<undefined | {[key:string]: any}>(undefined)
    useEffect(() => {
        values.current = selectedRecord?.values
    }, [selectedRecord])


    useEffect(() => {
        store.dispatch(fetchDatabases())
        store.dispatch(fetchForms())
    }, [])

    useEffect(() => {
        if (!selectedDatabase) {
            return
        }
        if (!selectedFormIntf) {
            return
        }
        store.dispatch(initNewRecord({
            databaseId: selectedDatabase.id,
            formId: selectedFormIntf.id,
            parentId: undefined
        }))
    }, [selectedDatabase?.id, selectedFormIntf?.id])

    if (!selectedDatabase) {
        return <Fragment/>
    }
    if (!selectedFormIntf) {
        return <Fragment/>
    }

    if (postRecordSuccess) {
        setTimeout(() => {
            store.dispatch(createRecordCancel())
        })
        if (!selectedForm){
            return <Fragment/>
        }
        return <Redirect to={`/databases/${selectedDatabase.id}/forms/${selectedForm.id}`}/>
    }

    const onSave = () => {
        if (selectedRecord) {
            store.dispatch(createRecord({record: selectedRecord}))
        }
    }

    const setValue = (key: string, value: any) => {
        if (selectedRecord) {
            store.dispatch(setRecordValue({id: selectedRecord?.id, key, value}))
        }
    }

    const openSubForm = (fieldId: string, subFormId: string) => {
        store.dispatch(recordOpenSubForm({id: subFormId}))
    }

    if (!values) {
        return <Fragment/>
    }

    return <RecordEditor
        onSave={() => onSave()}
        database={selectedDatabase}
        form={selectedFormIntf}
        onCancel={() => {
        }}
        parentForm={undefined}
        setValue={(key, value) => setValue(key, value)}
        success={postRecordSuccess}
        values={values}
        onOpenSubForm={({fieldId, subFormId}) => openSubForm(fieldId, subFormId)}
    />


}

type OpenSubFormFn = (params: { fieldId: string, subFormId: string }) => void

export type RecordEditorProps = {
    database: Database
    form: FormIntf
    parentForm: FormDefinition | undefined
    success: boolean
    setValue: (key: string, value: any) => void,
    values: { [key: string]: any }
    onSave: () => void
    onCancel: () => void
    onOpenSubForm: OpenSubFormFn
}

export const RecordEditor: FC<RecordEditorProps> = (props) => {

    const {
        database,
        form,
        values,
        setValue,
        onSave,
        onOpenSubForm
    } = props

    return (
        <Fragment>
            <Header
                form={form}
                database={database}
                title={`Add Record in ${form.name}`}/>
            <div className={" bg-dark text-light vh-100"}>
                <div className={"container py-4"}>
                    <div className={"row"}>
                        <div className={"col"}>
                            <h6 className={"text-uppercase"}>Add Record</h6>
                            <div className={"card bg-dark text-light border-secondary"}>
                                <div className={"card-body"}>
                                    {form?.fields?.map(f => {
                                        return <FormControl
                                            key={f.id}
                                            form={form}
                                            database={database}
                                            field={f}
                                            setValue={value => setValue(f.id, value)}
                                            value={values[f.id] ? values[f.id] : null}
                                            onOpenSubForm={onOpenSubForm}
                                        />
                                    })}
                                </div>
                            </div>
                            <div className={"d-flex flex-row p-3 "}>
                                <div className={"flex-grow-1"}/>

                                <button
                                    onClick={() => onSave()}
                                    className={"btn btn-lg btn-primary me-2"}>
                                    Save
                                </button>

                                <button className={"btn btn-lg btn-primary"}>Cancel</button>

                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </Fragment>
    )

}


type SetValueFn = (value: any) => void

type FormControlProps = {
    form: FormIntf
    database: Database
    field: FieldDefinition
    setValue: SetValueFn
    value: any
    onOpenSubForm: OpenSubFormFn
}

const FormControl: FC<FormControlProps> = props => {
    const {field} = props
    if (field.fieldType.text) {
        return TextFormControl(props)
    }
    if (field.fieldType.subForm) {
        return SubFormControl(props)
    }
    return <Fragment/>
}

const TextFormControl: FC<FormControlProps> = props => {
    return (
        <div className={"form-group mb-3"}>
            <label className={"form-label"}>{props.field.name}</label>
            <input
                onChange={(e) => props.setValue(e.target.value)}
                className={"form-control"}
                value={props.value ? props.value : ""}
                type={"text"}/>
        </div>
    )
}

const SubFormControl: FC<FormControlProps> = props => {

    const subForm = props.field.fieldType.subForm
    if (!subForm) {
        return <Fragment/>
    }

    return <div className={"mb-3"}>
        <h5>{props.field.name}</h5>
        <small
            className={"text-muted"}>{props.field.description}</small>
        <button
            onClick={() => props.onOpenSubForm({fieldId: props.field.id, subFormId: subForm.id})}
            className="btn w-100 btn-primary"
            type="button">
            Add record in {props.field.name}
        </button>
    </div>
}

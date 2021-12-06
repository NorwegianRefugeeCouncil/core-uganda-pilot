import React, {FC, Fragment, useCallback, useEffect, useState} from "react";
import {
    FormValue,
    postRecord,
    recorderActions, resetForm,
    selectCurrentForm,
    selectCurrentRecord,
    selectCurrentRootForm, selectPostRecords, selectSubRecords
} from "./recorder.slice";
import {useAppDispatch, useAppSelector} from "../../app/hooks";
import {FieldDefinition} from "core-js-api-client";
import {FieldEditor} from "./FieldEditor";
import {fetchDatabases} from "../../reducers/database";
import {fetchFolders} from "../../reducers/folder";
import {fetchForms, selectRootForm} from "../../reducers/form";
import {useParams, useLocation} from "react-router-dom"

export type RecordEditorProps = {
    formName: string
    fields: FieldDefinition[]
    values: { [key: string]: any }
    setValue: (key: string, value: any) => void
    selectSubRecord: (subRecordId: string) => void
    addSubRecord: (ownerFieldId: string) => void
    subRecords: {[key: string]: FormValue[]}
    saveRecord: () => void
}

export const RecordEditor: FC<RecordEditorProps> = props => {

    if (!props.fields) {
        return <Fragment/>
    }

    return <div className={"flex-grow-1 w-100 h-100 bg-dark text-light py-3 overflow-scroll"}>
        <div className={"container-fluid"}>
            <div className={"row justify-content-center"}>
                <div className={"col-12 col-lg-8"}>
                    <h4 className={"mb-4"}>Add record</h4>
                    <div className={"card bg-dark text-light border-secondary"}>
                        <div className={"card-body"}>
                            {props.fields.map(field => {
                                const value = props.values[field.id]
                                const setValue = (value: any) => {
                                    props.setValue(field.id, value)
                                }
                                const addSubRecord = () => {
                                    props.addSubRecord(field.id)
                                }
                                return <FieldEditor
                                    key={field.id}
                                    field={field}
                                    value={value}
                                    setValue={setValue}
                                    subRecords={props.subRecords[field.id]}
                                    selectSubRecord={props.selectSubRecord}
                                    addSubRecord={addSubRecord}/>
                            })}
                            <div className={"my-3"}>
                                <button
                                    onClick={() => props.saveRecord()}
                                    className={"btn btn-primary"}>
                                    Save Record
                                </button>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
}

export const RecordEditorContainer: FC<{}> = props => {

    const dispatch = useAppDispatch()

    // load data
    useEffect(() => {
        dispatch(fetchDatabases())
        dispatch(fetchFolders())
        dispatch(fetchForms())
    }, [dispatch])


    const params = useParams<{ formId: string }>()
    const location = useLocation()

    const [ownerRecordId, setOwnerRecordId] = useState<string | undefined>(undefined)
    const formIdFromPath = params.formId
    const currentRootForm = useAppSelector(selectCurrentRootForm)
    const rootFormFromPath = useAppSelector(s => selectRootForm(s, formIdFromPath))
    const currentForm = useAppSelector(selectCurrentForm)
    const currentRecord = useAppSelector(selectCurrentRecord)
    const subRecords = useAppSelector(state => {
        if (currentRecord){
            return selectSubRecords(state, currentRecord.recordId)
        }
        return {}
    })

    useEffect(() => {
        const search = new URLSearchParams(location.search)
        let ownerRecordIdFromQryParam = search.get("ownerRecordId");
        if (ownerRecordIdFromQryParam !== ownerRecordId) {
            setOwnerRecordId(ownerRecordIdFromQryParam ? ownerRecordIdFromQryParam : undefined)
        }
    }, [ownerRecordId, location])

    // make sure the form being edited is the one selected in the path
    useEffect(() => {

        if (!rootFormFromPath) {
            return
        }

        if (rootFormFromPath.id !== currentRootForm?.id) {
            dispatch(resetForm({
                formId: formIdFromPath,
                ownerId: ownerRecordId,
            }))
        }

    }, [dispatch, ownerRecordId, formIdFromPath, currentRootForm, rootFormFromPath])

    const setFieldValue = useCallback((key: string, value: any) => {
        if (currentRecord) {
            dispatch(recorderActions.setFieldValue({recordId: currentRecord.recordId, fieldId: key, value: value}))
        }
    }, [dispatch, currentRecord])

    const addSubRecord = useCallback((ownerFieldId: string) => {
        if (!currentRecord) {
            return
        }
        if (!currentForm) {
            return
        }
        const field = currentForm.fields.find(f => f.id === ownerFieldId)
        if (!field) {
            return
        }
        if (!field.fieldType.subForm) {
            return
        }
        const subFormId = field.id

        dispatch(recorderActions.addSubRecord({
            formId: subFormId,
            ownerFieldId: ownerFieldId,
            ownerRecordId: currentRecord.recordId
        }))

    }, [dispatch, currentForm, currentRecord])

    const recordsToPost = useAppSelector(selectPostRecords)

    const saveRecord = useCallback(() => {
        // do not save if we are not positioned on a record (should not happen)
        if (!currentRecord) {
            return
        }
        // do not save if we are not positioned on a form (should not happen)
        if (!currentForm) {
            return
        }

        if (currentRecord.formId !== formIdFromPath) {
            if (currentRecord.ownerRecordId) {
                dispatch(recorderActions.selectRecord({recordId: currentRecord.ownerRecordId}))
            }
        } else {
            dispatch(postRecord(recordsToPost))
        }

    }, [dispatch, formIdFromPath, currentRecord, recordsToPost, currentForm])

    const selectSubRecord = useCallback((subRecordId: string) => {
        dispatch(recorderActions.selectRecord({recordId: subRecordId}))
    }, [dispatch])

    if (!currentForm) {
        return <Fragment/>
    }

    if (!currentRecord) {
        return <Fragment/>
    }


    return <RecordEditor
        setValue={setFieldValue}
        fields={currentForm?.fields}
        values={currentRecord?.values}
        addSubRecord={addSubRecord}
        formName={currentForm.name}
        saveRecord={saveRecord}
        subRecords={subRecords}
        selectSubRecord={selectSubRecord}
    />
}

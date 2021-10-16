import React, {FC, Fragment, useCallback, useEffect, useState} from "react";
import {fetchForms, selectFormOrSubFormById, selectRootForm} from "../../reducers/form";
import {useAppDispatch, useAppSelector} from "../../app/hooks";
import {FieldDefinition, Record} from "../../types/types";
import {fetchDatabases} from "../../reducers/database";
import {fetchFolders} from "../../reducers/folder";
import {fetchRecords, selectRecordsForForm, selectRecordsSubFormCounts} from "../../reducers/records";
import {Link} from "react-router-dom";

export type FormBrowserProps = {
    formId: string
    fields: FieldDefinition[]
    records: Record[]
    getSubFormSum: (recordId: string, fieldId: string) => number
    parentRecordId: string | undefined
}

type subFormCountFn = (recordId: string, fieldId: string) => number

function mapHeaderField(field: FieldDefinition) {
    return <th key={field.id}
               style={{minWidth: field.name.length * 15}}>
        <small
            style={{fontSize: "0.75rem"}}
            className={"text-muted text-uppercase"}>{field.name}</small>
    </th>
}

function mapHeaderFields(fields: FieldDefinition[]) {
    return <tr>
        <th/>
        {fields.map(f => mapHeaderField(f))}
    </tr>
}

function mapRecordCell(field: FieldDefinition, record: Record, getSubFormCount: subFormCountFn) {
    if (field.fieldType.subForm) {
        const count = getSubFormCount(record.id, field.id)
        return <td key={field.id}>
            <span>
                <Link to={`/browse/forms/${field.fieldType.subForm.id}?parentRecordId=${record.id}`}>{count} records</Link>
            </span>
        </td>
    }

    return <td key={field.id}>
        <span className={"text-secondary"}>
        {record.values[field.id]}
        </span>
    </td>
}

function mapRecord(fields: FieldDefinition[], record: Record, getSubFormCount: subFormCountFn) {
    return <tr
        key={record.id}>
        <td><Link to={`/browse/records/${record.id}`}>
            <i className={"bi bi-search"}/>
        </Link></td>
        {fields.map(f => mapRecordCell(f, record, getSubFormCount))}
    </tr>
}

function mapRecords(fields: FieldDefinition[], records: Record[], getSubFormCount: subFormCountFn) {
    return records.map(r => mapRecord(fields, r, getSubFormCount))
}

export const FormBrowser: FC<FormBrowserProps> = props => {
    const {fields, records, formId, getSubFormSum, parentRecordId} = props

    let addRecordURL = `/edit/forms/${formId}/record`;
    if (parentRecordId) {
        addRecordURL += `?parentRecordId=${parentRecordId}`
    }

    return <div className={"flex-grow-1 w-100 h-100 overflow-scroll"}>
        <div className={"py-3 px-2"}>
            <Link className={"btn btn-primary"} to={addRecordURL}>Add Record</Link>
        </div>
        <table className={"table table-bordered w-100"}>
            <thead style={{lineHeight: "0.75rem"}}>
            {mapHeaderFields(fields)}
            </thead>
            <tbody style={{borderColor: "#dee2e6"}}>
            {mapRecords(fields, records, getSubFormSum)}
            </tbody>
        </table>
    </div>

}

export type FormBrowserContainerProps = {
    formId: string
    parentRecordId: string
}

export const FormBrowserContainer: FC<FormBrowserContainerProps> = props => {

    const dispatch = useAppDispatch()

    useEffect(() => {
        dispatch(fetchDatabases())
        dispatch(fetchFolders())
        dispatch(fetchForms())
    }, [dispatch])

    const form = useAppSelector((s) => {
        return selectFormOrSubFormById(s, props.formId)
    })

    const rootForm = useAppSelector(s => {
        return selectRootForm(s, props.formId)
    })

    const records = useAppSelector((s) => selectRecordsForForm(s, props.formId, props.parentRecordId))
    const subFormTotals = useAppSelector(selectRecordsSubFormCounts(form?.id))

    const getSubFormTotal = useCallback((recordId, fieldId) => {
        if (!subFormTotals.hasOwnProperty(recordId)) {
            return 0
        }
        if (!subFormTotals[recordId].hasOwnProperty(fieldId)) {
            return 0
        }
        return subFormTotals[recordId][fieldId]
    }, [subFormTotals])

    const [fetched, setFetched] = useState(false)

    useEffect(() => {
        if (!rootForm) {
            return
        }
        if (!form) {
            return
        }
        if (fetched) {
            return
        }
        setFetched(true)
        dispatch(fetchRecords({databaseId: rootForm.databaseId, formId: form?.id}))
        for (let field of form.fields) {
            if (field.fieldType.subForm) {
                dispatch(fetchRecords({databaseId: rootForm.databaseId, formId: field.fieldType.subForm.id}))
            }
        }
    }, [dispatch, rootForm, form, fetched])

    if (!form) {
        return <Fragment/>
    }


    return <FormBrowser
        getSubFormSum={getSubFormTotal}
        parentRecordId={props.parentRecordId}
        formId={form.id}
        fields={form?.fields}
        records={records}/>

}

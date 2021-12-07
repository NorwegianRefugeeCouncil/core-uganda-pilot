import React, {FC, Fragment, useCallback, useEffect, useState} from "react";
import {fetchForms, selectFormOrSubFormById, selectRootForm} from "../../reducers/form";
import {useAppDispatch, useAppSelector} from "../../app/hooks";
import {FieldDefinition, Record} from "core-js-api-client";
import {fetchDatabases} from "../../reducers/database";
import {fetchFolders} from "../../reducers/folder";
import {fetchRecords, selectRecordsForForm, selectRecordsSubFormCounts} from "../../reducers/records";
import {Link} from "react-router-dom";
import format from "date-fns/format"

type subFormCountFn = (recordId: string, fieldId: string) => number

type HeaderFieldProps = {
    field: FieldDefinition
    columnWidth: number
}

const HeaderField: FC<HeaderFieldProps> = props => {
    const {field, columnWidth} = props

    return <th key={field.id}
               className={"position-relative"}
               style={{minWidth: columnWidth, maxWidth: columnWidth}}>
        <div className={"d-flex flex-row align-items-center"}>
            <small
                style={{fontSize: "0.75rem"}}
                className={"text-muted text-uppercase"}>{field.name}</small>
        </div>
    </th>
}

type HeaderFieldsProps = {
    fields: FieldDefinition[]
    columnWidths: { [fieldId: string]: number }
}

export const HeaderFields: FC<HeaderFieldsProps> = props => {
    const {fields, columnWidths} = props
    return (
        <tr>
            <th style={{width: 35}}/>
            {fields.map(f => {
                return <HeaderField
                    key={f.id}
                    field={f}
                    columnWidth={columnWidths[f.id]}
                />
            })}
        </tr>
    )
}

function mapRecordCell(field: FieldDefinition, record: Record, getSubFormCount: subFormCountFn) {
    if (field.fieldType.subForm) {
        const count = getSubFormCount(record.id, field.id)
        return <td key={field.id}>
            <span>
                <Link
                    to={`/browse/forms/${field.id}?ownerRecordId=${record.id}`}>{count} records</Link>
            </span>
        </td>
    }

    const fieldValue = record.values.find((v: any) => v.fieldId === field.id)

    if (field.fieldType.month) {

        let date: Date | undefined
        if (fieldValue) {
            date = new Date(fieldValue.value)
        }

        return <td key={field.id} className={"text-secondary"}
                   style={{overflow: "hidden", textOverflow: "ellipsis", whiteSpace: "nowrap", wordBreak: "break-all"}}>
            {date ? format(date, "yyyy-MM") : ""}
        </td>
    }

    if (field.fieldType.week) {
        return <td key={field.id} className={"text-secondary"}
               style={{overflow: "hidden", textOverflow: "ellipsis", whiteSpace: "nowrap", wordBreak: "break-all"}}>
            {format(new Date(record.values[field.id]), "yyyy-'W'ww")}
        </td>
    }

    if (field.fieldType.reference) {
        return <td key={field.id}>
            <span>
                <Link to={`/browse/forms/${field.fieldType.reference.formId}`}>View</Link>
            </span>
        </td>
    }

    return <td key={field.id} className={"text-secondary"}
               style={{overflow: "hidden", textOverflow: "ellipsis", whiteSpace: "nowrap", wordBreak: "break-all"}}>
        {fieldValue?.value}
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

export type FormBrowserProps = {
    formId: string
    fields: FieldDefinition[]
    records: Record[]
    getSubFormSum: (recordId: string, fieldId: string) => number
    ownerRecordId: string | undefined
    columnWidths: { [fieldId: string]: number }
}

export const FormBrowser: FC<FormBrowserProps> = props => {
    const {fields, records, formId, getSubFormSum, ownerRecordId, columnWidths} = props

    let addRecordURL = `/edit/forms/${formId}/record`;
    if (ownerRecordId) {
        addRecordURL += `?ownerRecordId=${ownerRecordId}`
    }

    return <div className={"flex-grow-1 w-100 h-100 overflow-scroll bg-light"}>
        <div className={"py-3 px-2"}>
            <Link className={"btn btn-primary"} to={addRecordURL}>Add Record</Link>
        </div>
        <div className={"px-2"}>
            <table className={"table shadow bg-white table-bordered w-100"} style={{tableLayout: "fixed"}}>
                <thead style={{lineHeight: "0.75rem"}}>
                <HeaderFields fields={fields} columnWidths={columnWidths}/>
                </thead>
                <tbody style={{borderColor: "#dee2e6", borderTop: "none"}}>
                {mapRecords(fields, records, getSubFormSum)}
                </tbody>
            </table>
        </div>
    </div>

}

export type FormBrowserContainerProps = {
    formId: string
    ownerRecordId: string
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

    const records = useAppSelector((s) => selectRecordsForForm(s, props.formId, props.ownerRecordId))
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
                dispatch(fetchRecords({databaseId: rootForm.databaseId, formId: field.id}))
            }
        }
    }, [dispatch, rootForm, form, fetched])


    const [columnWidths, setColumnWidths] = useState<{ [key: string]: number }>({})

    useEffect(() => {
        if (!form?.fields) {
            return
        }
        const widths: { [key: string]: number } = {}
        for (let field of form?.fields) {
            widths[field.id] = 200
        }
        setColumnWidths(widths)
    }, [form?.fields])

    if (!form) {
        return <Fragment/>
    }

    return <FormBrowser
        getSubFormSum={getSubFormTotal}
        ownerRecordId={props.ownerRecordId}
        formId={form.id}
        fields={form.fields}
        records={records}
        columnWidths={columnWidths}
    />

}

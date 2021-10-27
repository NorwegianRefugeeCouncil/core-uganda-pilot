import React, {FC, useCallback, useEffect, useState} from "react";
import {Record} from "../types/types";
import {useAppDispatch, useAppSelector} from "../app/hooks";
import {fetchRecords, recordGlobalSelectors, selectRecords} from "../reducers/records";
import {useDispatch} from "react-redux";
import {formGlobalSelectors, selectFormOrSubFormById, selectRootForm} from "../reducers/form";

export type RecordPickerProps = {
    disabled?: boolean
    recordId: string | undefined
    setRecordId: (recordId: string | undefined) => void
    records: Record[]
}

export const RecordPicker: FC<RecordPickerProps> = props => {
    const {
        disabled,
        recordId,
        setRecordId,
        records
    } = props

    return <div>
        <select
            disabled={disabled}
            onChange={e => setRecordId(e.target.value)}
            value={recordId ? recordId : ""}
            className="form-select"
            aria-label="Select Record">
            <option disabled={true} value={""}>{"No Records"}</option>
            {records.map(r => {
                return (
                    <option
                        value={r.id}>{r.id}
                    </option>
                );
            })}
        </select>
    </div>

}

export type RecordPickerContainerProps = {
    recordId: string | undefined
    setRecordId?: (recordId: string | undefined) => void
    setRecord?: (record: Record | undefined) => void
    parentId?: string
    formId?: string
}

export const RecordPickerContainer: FC<RecordPickerContainerProps> = props => {

    const {
        parentId,
        formId,
        recordId,
        setRecordId,
        setRecord
    } = props

    const dispatch = useAppDispatch()

    const form = useAppSelector(state => {
        return selectFormOrSubFormById(state, props.formId ? props.formId : "")
    })

    const rootForm = useAppSelector(state => {
        return selectRootForm(state, form ? form.id : "")
    })

    const [pending, setPending] = useState(false)

    useEffect(() => {
        if (rootForm && form && !pending) {
            dispatch(fetchRecords({databaseId: rootForm.databaseId, formId: form?.id})).then(() => {
                setPending(true)
            }).catch(() => {
                setPending(false)
            })
        }
    }, [dispatch, form, pending, rootForm])

    const records = useAppSelector(state => {
        return selectRecords(state, {parentId, formId})
    })

    const record = useAppSelector(state => {
        return recordGlobalSelectors.selectById(state, recordId ? recordId : "")
    })

    const callback = useCallback((recordId: string | undefined) => {
        if (setRecord) {
            setRecord(record)
        }
        if (setRecordId) {
            setRecordId(recordId)
        }
    }, [record, setRecord, setRecordId])

    return <RecordPicker
        disabled={false}
        recordId={recordId}
        setRecordId={callback}
        records={records}/>

}

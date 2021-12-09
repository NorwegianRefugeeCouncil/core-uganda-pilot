import React, {FC, useCallback, useEffect, useState} from "react";
import {Record} from "core-api-client";
import {useAppDispatch, useAppSelector} from "../app/hooks";
import {fetchRecords, recordGlobalSelectors, selectRecords} from "../reducers/records";
import {selectFormOrSubFormById, selectRootForm} from "../reducers/form";

export type RecordPickerProps = {
    disabled?: boolean
    recordId: string | undefined
    setRecordId: (recordId: string | undefined) => void
    records: Record[]
    getDisplayStr: (record: Record) => string
}

export const RecordPicker: FC<RecordPickerProps> = props => {
    const {
        disabled,
        recordId,
        setRecordId,
        records,
        getDisplayStr
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
                        value={r.id}>{getDisplayStr(r)}
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
    ownerId?: string
    formId?: string
}

export const RecordPickerContainer: FC<RecordPickerContainerProps> = props => {

    const {
        ownerId,
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
        return selectRecords(state, {ownerId: ownerId, formId})
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
        records={records}
        getDisplayStr={r => {
            let result = ""
            if (!form) {
                return result
            }
            for (const field of form?.fields) {
                if (field.key) {
                    result += r.values.find((v: any) => v.fieldId === field.id)
                }
            }
            return result
        }
        }
    />

}

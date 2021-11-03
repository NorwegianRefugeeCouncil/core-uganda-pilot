import {FieldDefinition} from "core-js-api-client";
import React, {FC, Fragment, useState} from "react";
import {FormValue} from "./recorder.slice";
import {RecordPickerContainer} from "../../components/RecordPicker";
import format from "date-fns/format"

export type FieldEditorProps = {
    field: FieldDefinition
    value: any
    setValue: (value: any) => void
    addSubRecord: () => void
    selectSubRecord: (subRecordId: string) => void
    subRecords: FormValue[] | undefined
}

export const ReferenceFieldEditor: FC<FieldEditorProps> = props => {
    const {field, value, setValue} = props
    return <div className={"form-group mb-2"}>
        <label
            className={"form-label opacity-75"}
            htmlFor={field.id}>{field.name}</label>
        <RecordPickerContainer
            formId={field.fieldType.reference?.formId}
            recordId={value}
            setRecordId={setValue}/>
        {mapFieldDescription(field)}
    </div>
}


export const TextFieldEditor: FC<FieldEditorProps> = props => {
    const {field, value, setValue} = props
    return <div className={"form-group mb-2"}>
        <label
            className={"form-label opacity-75"}
            htmlFor={field.id}>{field.name}</label>
        <input
            className={"form-control bg-dark text-light border-secondary"}
            type={"text"}
            id={field.id} value={value ? value : ""}
            onChange={event => setValue(event.target.value)}/>
        {mapFieldDescription(field)}
    </div>
}

export const MultilineTextFieldEditor: FC<FieldEditorProps> = props => {
    const {field, value, setValue} = props
    return <div className={"form-group mb-2"}>
        <label
            className={"form-label opacity-75"}
            htmlFor={field.id}>{field.name}</label>
        <textarea
            className={"form-control bg-dark text-light border-secondary"}
            id={field.id} value={value ? value : ""}
            onChange={event => setValue(event.target.value)}/>
        {mapFieldDescription(field)}
    </div>
}

export const DateFieldEditor: FC<FieldEditorProps> = props => {
    const {field, value, setValue} = props
    return <div className={"form-group mb-2"}>
        <label
            className={"form-label opacity-75"}
            htmlFor={field.id}>{field.name}</label>
        <input
            className={"form-control bg-dark text-light border-secondary"}
            type={"date"}
            id={field.id} value={value ? value : ""}
            onChange={event => setValue(event.target.value)}/>
        {mapFieldDescription(field)}
    </div>
}

export const MonthFieldEditor: FC<FieldEditorProps> = props => {
    const {field, value, setValue} = props
    const expectedLength = 7;

    const [localValue, setLocalValue] = useState<string>(() => {
        if (!value) {
            return ""
        }
        try {
            return format(value, "yyyy-MM")
        } catch (e) {
            return ""
        }
    })

    const isValidLength = () => localValue.length === expectedLength;

    function isValid(s: string) {
        const valid = /^(?:19|20|21)\d{2}-[01]\d$/
        const m = +s.slice(5)
        return valid.test(s) && 0 < m && m <= 12;
    }

    return <div className={"form-group mb-2"}>
        <label
            className={"form-label opacity-75"}
            htmlFor={field.id}>{field.name}</label>
        <input
            className={`form-control bg-dark text-light border-secondary ${!isValid(localValue) && isValidLength() ? " is-invalid" : ""}`}
            type={"month"}
            maxLength={expectedLength}
            id={field.id} value={localValue ? localValue : ""}
            pattern={"[0-9]{4}-[0-9]{2}"}
            placeholder={"YYYY-MM"}
            onChange={event => {
                const v = event.target.value;
                setLocalValue(v);
                if (!isValid(v)) return
                const date = new Date(+v.slice(0, 4), +v.slice(5, 7) - 1, 1)
                setValue(date);
            }}
        />

        {mapFieldDescription(field)}
    </div>
}

export const WeekFieldEditor: FC<FieldEditorProps> = props => {
    const {field, value, setValue} = props
    return <div className={"form-group mb-2"}>
        <label
            className={"form-label opacity-75"}
            htmlFor={field.id}>{field.name}</label>
        <input
            className={"form-control bg-dark text-light border-secondary"}
            type={"week"}
            id={field.id} value={value ? value : ""}
            onChange={event => setValue(event.target.value)}/>
        {mapFieldDescription(field)}
    </div>
}

export const QuantityFieldEditor: FC<FieldEditorProps> = props => {
    const {field, value, setValue} = props
    return <div className={"form-group mb-2"}>
        <label
            className={"form-label opacity-75"}
            htmlFor={field.id}>{field.name}</label>
        <input
            className={"form-control bg-dark text-light border-secondary"}
            type="number"
            id={field.id} value={value ? value : ""}
            onChange={event => setValue(event.target.value)}/>
        {mapFieldDescription(field)}
    </div>
}

export const SingleSelectFieldEditor: FC<FieldEditorProps> = props => {
    const {field, value, setValue} = props
    return <div className={"form-group mb-2"}>
        <label
            className={"form-label opacity-75"}
            htmlFor={field.id}>{field.name}</label>
        <select
            className={"form-control bg-dark text-light border-secondary"}
            id={field.id} value={value ? value : ""}
            onChange={event => setValue(event.target.value)}>
            <option disabled={field.required || field.key} value={""}/>
            {/** TODO field.options.map(o => <option value={o}>{o}</option>)**/}
        </select>
        {mapFieldDescription(field)}
    </div>
}

function subRecord(record: FormValue, select: () => void) {
    return <a href="/#" key={record.recordId}
              onClick={(e) => {
                  e.preventDefault()
                  select()
              }}
              className={"list-group-item list-group-item-action bg-dark border-secondary text-secondary"}>
        View Record
    </a>
}

function subRecords(records: FormValue[], select: (id: string) => void) {
    return <div className={"list-group bg-dark mb-3"}>
        {records.map(r => subRecord(r, () => {
            select(r.recordId)
        }))}
    </div>
}


export const SubFormFieldEditor: FC<FieldEditorProps> = props => {
    const {field, addSubRecord} = props
    return <div className={"mb-2"}>
        <div className={"bg-primary border-2"}/>
        <label className={"form-label opacity-75"}>{field.name}</label>
        {props.subRecords ? subRecords(props.subRecords, props.selectSubRecord) : <Fragment/>}
        <button
            onClick={addSubRecord}
            className={"btn btn-sm btn-outline-primary w-100"}>Add record in {field.name}</button>
        {mapFieldDescription(field)}
    </div>
}


export const FieldEditor: FC<FieldEditorProps> = props => {
    const {fieldType} = props.field
    if (fieldType.text) {
        return <TextFieldEditor {...props} />
    } else if (fieldType.week) {
        return <WeekFieldEditor {...props} />
    } else if (fieldType.subForm) {
        return <SubFormFieldEditor {...props} />
    } else if (fieldType.reference) {
        return <ReferenceFieldEditor {...props} />
    } else if (fieldType.multilineText) {
        return <MultilineTextFieldEditor {...props} />
    } else if (fieldType.date) {
        return <DateFieldEditor {...props} />
    } else if (fieldType.month) {
        return <MonthFieldEditor {...props} />
    } else if (fieldType.quantity) {
        return <QuantityFieldEditor {...props} />
    } else if (fieldType.singleSelect) {
        return <SingleSelectFieldEditor {...props} />
    } else {
        return <Fragment/>
    }
}

export const mapFieldDescription = (fd: FieldDefinition) => {
    if (fd.description) {
        return <small className={"text-muted"}>{fd.description}</small>
    } else {
        return <Fragment/>
    }
}

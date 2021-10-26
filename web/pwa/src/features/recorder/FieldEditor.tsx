import {FieldDefinition} from "../../types/types";
import React, {FC, Fragment} from "react";
import {FormValue} from "./recorder.slice";
import {RecordPickerContainer} from "../../components/RecordPicker";

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

export const MultilineFieldEditor: FC<FieldEditorProps> = props => {
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


function subRecord(record: FormValue, select: () => void) {
    return <a href={"#"} key={record.recordId}
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
    console.log(props.subRecords)
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
    } else if (fieldType.subForm) {
        return <SubFormFieldEditor {...props} />
    } else if (fieldType.reference) {
        return <ReferenceFieldEditor {...props} />
    } else if (fieldType.multiline) {
        return <MultilineFieldEditor {...props} />
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

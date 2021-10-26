import {FieldKind} from "../../types/types";
import React, {FC, Fragment} from "react";
import DatabasePickerContainer from "../../components/DatabasePicker";
import FormPickerContainer from "../../components/FormPicker";

type FormerFieldProps = {
    isSelected: boolean
    selectField: () => void
    fieldType: FieldKind
    fieldName: string
    setFieldName: (fieldName: string) => void
    fieldRequired: boolean,
    setFieldRequired: (required: boolean) => void
    fieldIsKey: boolean,
    setFieldIsKey: (isKey: boolean) => void
    fieldDescription: string,
    setFieldDescription: (description: string) => void
    openSubForm: () => void
    saveField: () => void
    setReferencedDatabaseId: (databaseId: string) => void
    referencedDatabaseId: string | undefined
    setReferencedFormId: (formId: string) => void
    referencedFormId: string | undefined
    cancel: () => void
}
export const FormerField: FC<FormerFieldProps> = props => {

    const {
        fieldName,
        setFieldName,
        isSelected,
        selectField,
        fieldDescription,
        setFieldDescription,
        fieldType,
        referencedDatabaseId,
        referencedFormId,
        setReferencedDatabaseId,
        setReferencedFormId,
        openSubForm,
        saveField,
        cancel
    } = props

    if (!isSelected) {
        return <div
            onClick={() => selectField()}
            style={{cursor: "pointer"}}
            className={"card bg-dark text-light border-light mb-2"}>
            <div className={"card-body p-3"}>
                <div className={"d-flex flex-row"}>
                    <span className={"flex-grow-1"}>{fieldName}</span>
                    <small className={"text-uppercase"}>{fieldType}</small>
                </div>

            </div>
        </div>
    }

    return <div className={"card text-dark"}>
        <div className={"card-body"}>
            <h6 className={"card-title text-uppercase"}>{fieldType}</h6>
            <div className={"form-group mb-2"}>
                <label className={"form-label"} htmlFor={"fieldName"}>Field Name</label>
                <input className={"form-control"}
                       id={"database"}
                       type={"text"}
                       value={fieldName ? fieldName : ""}
                       onChange={event => setFieldName(event.target.value)}/>
            </div>
            <div className={"form-group mb-2"}>
                <label className={"form-label"} htmlFor={"fieldName"}>Description</label>
                <textarea
                    className={"form-control"}
                    id={"fieldDescription"}
                    onChange={(event) => setFieldDescription(event.target.value)}
                    value={fieldDescription ? fieldDescription : ""}
                />
            </div>
            {fieldType === "subform"
                ? <button className={"btn btn-primary"} onClick={() => openSubForm()}>Open Sub Form</button>
                : <Fragment/>
            }
            {fieldType === FieldKind.Reference
                ? (
                    <div>
                        <div className={"form-group mb-2"}>
                            <label className={"form-label"}>Database</label>
                            <DatabasePickerContainer
                                setDatabaseId={setReferencedDatabaseId}
                                databaseId={referencedDatabaseId}/>
                        </div>
                        <div className={"form-group"}>
                            <label className={"form-label"}>Form</label>
                            <FormPickerContainer databaseId={referencedDatabaseId}
                                                 formId={referencedFormId}
                                                 setFormId={setReferencedFormId}/>
                        </div>
                    </div>
                ) : <Fragment/>}
        </div>
        <div className={"card-footer"}>
            <button onClick={() => saveField()} className={"btn btn-primary"}>Save</button>
            <button onClick={() => cancel()} className={"btn btn-secondary"}>Cancel</button>
        </div>
    </div>
}

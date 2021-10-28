import {FieldKind} from "../../types/types";
import React, {FC, Fragment} from "react";
import DatabasePickerContainer from "../../components/DatabasePicker";
import FormPickerContainer from "../../components/FormPicker";

type FormerFieldProps = {
    isSelected: boolean
    selectField: () => void
    fieldType: FieldKind
    fieldOptions?: string[]
    setFieldOption: (i: number, value: string) => void
    addOption: () => void
    removeOption: (index: number) => void
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
        fieldOptions,
        setFieldOption,
        addOption,
        removeOption,
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
        cancel,
        fieldIsKey,
        setFieldIsKey,
        fieldRequired,
        setFieldRequired,
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


            {/* Form Title */}

            <h6 className={"card-title text-uppercase"}>{fieldType}</h6>


            <div className={"row"}>

                {/* Left Hand Side Section */}

                <div className={"col-8"}>


                    {/* Form Name */}

                    <div className={"form-group mb-2"}>
                        <label className={"form-label"} htmlFor={"fieldName"}>Field Name</label>
                        <input className={"form-control"}
                               id={"database"}
                               type={"text"}
                               value={fieldName ? fieldName : ""}
                               onChange={event => setFieldName(event.target.value)}/>
                    </div>

                    {/* Options */}

                    {fieldType === "singleSelect" &&
                    (
                        <div className={"form-group mb-2"}>
                            <div className="d-flex justify-content-between align-items-center mb-2">
                                <label className={"form-label"} htmlFor={"fieldName"}>Field Options</label>
                                <button type="button" className="btn btn-outline-primary"
                                        onClick={() => addOption()}>Add option
                                </button>
                            </div>
                            {fieldOptions?.map((opt, i) => (
                                <div key={i} className="d-flex mb-2">
                                    <input className="form-control me-3"
                                           id={`fieldOption-${i}`}
                                           type={"text"}
                                           value={opt ? opt : ""}
                                           onChange={event => setFieldOption(i, event.target.value)}/>
                                    <button type="button" className="btn btn-outline-danger"
                                            onClick={() => removeOption(i)}><i className="bi bi-x"/>
                                    </button>
                                </div>
                            ))}
                        </div>
                    )
                    }

                    {/* Form Description */}

                    <div className={"form-group mb-2"}>
                        <label className={"form-label"} htmlFor={"fieldName"}>Description</label>
                        <textarea
                            className={"form-control"}
                            id={"fieldDescription"}
                            onChange={(event) => setFieldDescription(event.target.value)}
                            value={fieldDescription ? fieldDescription : ""}
                        />
                    </div>


                    {/* Open Subform Button */}

                    {fieldType === "subform"
                        ? <button className={"btn btn-primary"} onClick={() => openSubForm()}>Open Sub Form</button>
                        : <Fragment/>
                    }


                    {/* Configure Reference Field */}

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


                {/* Right Hand Side Section */}

                <div className={"col-4"}>


                    {/* Required Checkbox */}

                    <div className="form-check">
                        <input
                            disabled={fieldIsKey}
                            className="form-check-input"
                            type="checkbox"
                            value=""
                            onChange={() => setFieldRequired(!fieldRequired)}
                            checked={fieldRequired}
                            id="required"/>
                        <label className="form-check-label" htmlFor="required">
                            Required
                        </label>
                    </div>


                    {/* Key Checkbox */}

                    <div className="form-check">
                        <input
                            className="form-check-input"
                            type="checkbox"
                            value=""
                            onChange={() => setFieldIsKey(!fieldIsKey)}
                            checked={fieldIsKey}
                            id="key"/>
                        <label className="form-check-label" htmlFor="key">
                            Key
                        </label>
                    </div>

                </div>

            </div>

        </div>
        <div className={"card-footer"}>
            <button onClick={() => saveField()} className={"btn btn-primary me-2 shadow"}>Save</button>
            <button onClick={() => cancel()} className={"btn btn-secondary shadow"}>Cancel</button>
        </div>
    </div>
}

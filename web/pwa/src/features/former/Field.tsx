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
                            onClick={() => setFieldRequired(!fieldRequired)}
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
                            onClick={() => setFieldIsKey(!fieldIsKey)}
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

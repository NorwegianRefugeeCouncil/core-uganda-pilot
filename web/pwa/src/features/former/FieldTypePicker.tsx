import React, {FC} from "react";
import {FieldKind} from "../../types/types";

export type FieldTypePickerProps = {
    onSubmit: (fieldKind: FieldKind) => void
    onCancel: () => void
}

export const FieldTypePicker: FC<FieldTypePickerProps> = props => {
    const {onSubmit, onCancel} = props
    return <div className={"card bg-dark text-light border-secondary"}>
        <div className={"card-body bg-primary"}>
            <div className={"d-flex flex-row"}>

                <button
                    className={"btn btn-primary m-2 border-light"}
                    onClick={() => onSubmit(FieldKind.Text)}>
                    Text
                </button>                <button
                    className={"btn btn-primary m-2 border-light"}
                    onClick={() => onSubmit(FieldKind.MultilineText)}>
                    MultilineText Text
                </button>
                <button
                    className={"btn btn-primary m-2 border-light"}
                    onClick={() => onSubmit(FieldKind.SubForm)}>
                    Subform
                </button>
                <button
                    className={"btn btn-primary m-2 border-light"}
                    onClick={() => onSubmit(FieldKind.Reference)}>
                    Reference
                </button>
                <button
                    className={"btn btn-primary m-2 border-light"}
                    onClick={() => onSubmit(FieldKind.Date)}>
                    Date
                </button>
                <button
                    className={"btn btn-primary m-2 border-light"}
                    onClick={() => onSubmit(FieldKind.Quantity)}>
                    Quantity
                </button>

            </div>
        </div>
        <div className={"card-footer border-secondary"}>
            <button
                className={"btn btn-secondary m-2"}
                onClick={() => onCancel()}>
                Cancel
            </button>
        </div>
    </div>
}

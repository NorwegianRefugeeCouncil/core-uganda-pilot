import {FC, Fragment} from "react";
import {FieldKind} from "./client";
import {AddFieldProps, FieldPropsState} from "./store";


type FieldEditorProps = {
    onSave: () => void
    onOpenSubForm: () => void
    setKind: (kind: FieldKind) => void
    setName: (name: string) => void
    setCode: (code: string) => void
    setIsKey: (key: boolean) => void
    setIsRequired: (required: boolean) => void
    setDescription: (description: string) => void
} & FieldPropsState

export const FieldEditor: FC<FieldEditorProps> = (props) => {
    const {
        type,
        setKind,
        name,
        setName,
        code,
        setCode,
        isRequired,
        setIsRequired,
        isKey,
        setIsKey,
        description,
        setDescription,
        onSave,
        onOpenSubForm,
    } = props

    if (!type) {
        return GetFieldKindChooser({setFieldKind: setKind})
    }

    const getSubFormOpener = () => {
        if (type === FieldKind.SubForm) {
            return <div>
                <button className={"btn btn-primary mt-2"} onClick={() => onOpenSubForm()}>Open Sub Form</button>
            </div>
        }
    }

    return <Fragment>
        <div>

            <div className={"bg-light text-dark px-5 py-3"}>
                <small
                    className={"text-uppercase"}>{GetFieldKindIcon({fieldKind: type})} {GetFieldKindText(type)}
                </small>
            </div>

            <div className={"d-flex flex-row"}>
                <div className={"col-10 bg-white text-dark px-5 pt-3 py-2 pb-3"}>
                    <div className={"form-group"}>
                        <label htmlFor={"name"}>Field Name:</label>
                        <input
                            id={"name"}
                            name={"name"}
                            className={"form-control"}
                            type={"text"}
                            value={name}
                            onChange={e => setName(e.target.value)}/>
                    </div>

                    <div className={"form-group mt-2"}>
                        <label htmlFor={"name"}>Description:</label>
                        <textarea
                            id={"description"}
                            name={"description"}
                            className={"form-control"}
                            rows={7}
                            value={description}
                            onChange={e => setDescription(e.target.value)}/>
                    </div>

                    {getSubFormOpener()}


                </div>
                <div className={"col-2 bg-white text-dark pe-5"}>
                    <div className={"form-group"}>
                        <label htmlFor={"code"}>Code</label>
                        <input
                            id={"code"}
                            name={"code"}
                            className={"form-control"}
                            type={"text"}
                            value={code}
                            onChange={e => setCode(e.target.value)}/>
                    </div>

                    <div className="form-check mt-2">

                        <input
                            className="form-check-input"
                            type="checkbox"
                            checked={isKey}
                            onClick={() => setIsKey(!isKey)}
                            id="key"/>

                        <label className="form-check-label" htmlFor="key">
                            Key
                        </label>

                    </div>

                    <div className="form-check mt-2">

                        <input
                            className="form-check-input"
                            type="checkbox"
                            checked={isRequired}
                            onClick={() => setIsRequired(!isRequired)}
                            id="required"/>

                        <label className="form-check-label" htmlFor="required">
                            Required
                        </label>

                    </div>

                </div>
            </div>
            <div className={"col bg-light text-dark px-5 py-4"}>
                <div className={"d-flex flex-row"}>
                    <div className={"flex-grow-1"}/>

                    <button
                        className={"btn btn-primary"}
                        disabled={name.length === 0}
                        onClick={() => onSave()}>
                        Save
                    </button>

                    <button className={"btn btn-primary ms-2"}>
                        Delete
                    </button>

                </div>
            </div>
        </div>
    </Fragment>

}

type FieldKindSetter = (fieldKind: FieldKind) => void

function GetFieldKindChooser({setFieldKind}: { setFieldKind: FieldKindSetter }) {
    return <Fragment>
        <div className={"bg-light text-dark p-3"}>
            <h5>Select field type:</h5>

            <div className={"row"}>

                <div className={"col-4"}>
                    <button className={"btn btn-block btn-primary"} onClick={() => setFieldKind(FieldKind.Text)}>Text
                    </button>
                </div>

                <div className={"col-4"}>
                    <button className={"btn btn-block btn-primary"} onClick={() => setFieldKind(FieldKind.SubForm)}>Sub
                        Form
                    </button>
                </div>

            </div>

        </div>

    </Fragment>
}

export function GetFieldKindText(fieldKind: FieldKind | undefined): string {
    switch (fieldKind) {
        case FieldKind.Reference:
            return "Reference"
        case FieldKind.SubForm:
            return "Sub Form"
        case FieldKind.Text:
            return "Text"
    }
    return ""
}

const GetFieldKindIcon: FC<{ fieldKind: FieldKind | undefined }> = (props: { fieldKind: FieldKind | undefined }) => {
    const {fieldKind} = props
    switch (fieldKind) {
        case FieldKind.Reference:
            return <i className={"bi"}/>
        case FieldKind.SubForm:
            return <i className={"bi bi-diagram-2"}/>
        case FieldKind.Text:
            return <i className={"bi bi-type"}/>
    }
    return <Fragment/>
}

function getFieldKindForField(fieldDefinition: AddFieldProps): FieldKind | undefined {
    if (fieldDefinition.type === FieldKind.Text) {
        return FieldKind.Text
    }
    if (fieldDefinition.type === FieldKind.SubForm) {
        return FieldKind.SubForm
    }
    if (fieldDefinition.type === FieldKind.Reference) {
        return FieldKind.Reference
    }
    return undefined
}

export const GetFieldIcon: FC<{ fieldDef: AddFieldProps }> = (props) => {
    const {fieldDef} = props
    const fieldKind = getFieldKindForField(fieldDef)
    return GetFieldKindIcon({fieldKind})
}

export const GetFieldText = (props: { fieldDef: AddFieldProps }) => {
    const {fieldDef} = props
    const fieldKind = getFieldKindForField(fieldDef)
    return GetFieldKindText(fieldKind)
}

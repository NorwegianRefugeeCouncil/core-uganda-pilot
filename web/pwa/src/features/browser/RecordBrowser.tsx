import {Fragment, FC} from "react";
import {useFormOrSubForm, useParentRecord, useRecordFromPath, useSubRecords} from "../../app/hooks";
import {FieldDefinition, Record} from "../../types/types";
import {Link} from "react-router-dom";

type RecordFieldProps = {
    field: FieldDefinition
    value: any
    subRecords: Record[] | undefined
}
const RecordField: FC<RecordFieldProps> = props => {
    const {field, value} = props
    return <div className={"form-group mb-3"}>
        <label className={"form-label"}>{field.name}</label>
        {field.fieldType.reference
            ? (<div><Link to={`/browse/records/${value}`}>View</Link></div>)
            : <div className={"fw-bold"}>{value}</div>}
        {props.subRecords?.map(r => <Link to={`/browse/records/${r.id}`}>Sub Record</Link>)}
    </div>
}


function mapRecordField(record: Record, field: FieldDefinition, subRecords: Record[] | undefined) {
    const value = record.values[field.id]
    return <RecordField field={field} value={value} subRecords={subRecords}/>
}

export const RecordBrowser: FC = props => {
    const record = useRecordFromPath("recordId")
    const form = useFormOrSubForm(record?.formId)
    const subRecords = useSubRecords(record?.id)
    const parentRecord = useParentRecord(record?.id)

    if (!record) {
        return <div>Record not found</div>
    }
    if (!form) {
        return <div>Form not found</div>
    }
    return <div className={"flex-grow-1 bg-light py-3"}>
        <div className={"container"}>
            <div className={"row"}>
                <div className={"col"}>
                    <div className={"card shadow"}>
                        <div className={"card-body"}>
                            {parentRecord
                                ? <div className={"mb-2"}>
                                    <Link to={`/browse/records/${parentRecord.id}`}>Back to Parent
                                        Record
                                    </Link>
                                </div>
                                : <Fragment/>}
                            {form?.fields.map(f => mapRecordField(record, f, subRecords?.byFieldId[f.id]))}
                        </div>
                    </div>

                </div>
            </div>
        </div>

    </div>
}

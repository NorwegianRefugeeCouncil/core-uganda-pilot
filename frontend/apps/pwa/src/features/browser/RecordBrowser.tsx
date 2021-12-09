import { FC, Fragment } from 'react';
import { FieldDefinition, Record } from 'core-api-client';
import { Link } from 'react-router-dom';
import format from 'date-fns/format';

import { useFormOrSubForm, useOwnerRecord, useRecordFromPath, useSubRecords } from '../../app/hooks';

type RecordFieldProps = {
  field: FieldDefinition;
  value: any;
  subRecords: Record[] | undefined;
};
const RecordField: FC<RecordFieldProps> = (props) => {
  const { field, value } = props;

  function renderField(f: FieldDefinition, v: any) {
    if (f.fieldType.reference) {
      return (
        <div>
          <Link to={`/browse/records/${v}`}>View</Link>
        </div>
      );
    }
    if (f.fieldType.month) {
      return <div className="fw-bold">{format(new Date(v), 'yyyy-MM')}</div>;
    }
    if (f.fieldType.week) {
      return <div className="fw-bold">{format(new Date(v), "yyyy-'W'ww")}</div>;
    }
    return <div className="fw-bold">{v}</div>;
  }

  return (
    <div className="form-group mb-3">
      <label className="form-label">{field.name}</label>
      {renderField(field, value)}
      {props.subRecords?.map((r) => (
        <Link to={`/browse/records/${r.id}`}>Sub Record</Link>
      ))}
    </div>
  );
};

function mapRecordField(record: Record, field: FieldDefinition, subRecords: Record[] | undefined) {
  let value = '';
  const fieldValue = record.values.find((v: any) => v.fieldId === field.id);
  if (fieldValue) {
    value = fieldValue.value;
  }
  return <RecordField field={field} value={`${value}`} subRecords={subRecords} />;
}

export const RecordBrowser: FC = (props) => {
  const record = useRecordFromPath('recordId');
  const form = useFormOrSubForm(record?.formId);
  const subRecords = useSubRecords(record?.id);
  const ownerRecord = useOwnerRecord(record?.id);

  if (!record) {
    return <div>Record not found</div>;
  }
  if (!form) {
    return <div>Form not found</div>;
  }
  return (
    <div className="flex-grow-1 bg-light py-3">
      <div className="container">
        <div className="row">
          <div className="col">
            <div className="card shadow">
              <div className="card-body">
                {ownerRecord ? (
                  <div className="mb-2">
                    <Link to={`/browse/records/${ownerRecord.id}`}>Back to Parent Record</Link>
                  </div>
                ) : (
                  <></>
                )}
                {form?.fields.map((f) => mapRecordField(record, f, subRecords?.byFieldId[f.id]))}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

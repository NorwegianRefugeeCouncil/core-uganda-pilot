import { FC } from 'react';
import { FieldDefinition, Record, SelectOption } from 'core-api-client';
import { Link } from 'react-router-dom';

import {
  useFormOrSubForm,
  useOwnerRecord,
  useRecordFromPath,
  useSubRecords,
} from '../../app/hooks';

import { RecordField } from './RecordField';

function mapRecordField(
  record: Record,
  field: FieldDefinition,
  subRecords: Record[] | undefined,
) {
  let value = '';
  const fieldValue = record.values.find((v: any) => v.fieldId === field.id);
  if (fieldValue && typeof fieldValue.value === 'string') {
    value = fieldValue.value;
  }
  if (fieldValue && field.fieldType.multiSelect) {
    const selected = field.fieldType.multiSelect.options.filter(
      (o: SelectOption) => {
        if (fieldValue?.value == null) {
          return false;
        }
        return fieldValue.value.includes(o.id);
      },
    );

    value = selected.map((s) => s.name).join(', ');
  }
  return (
    <RecordField
      key={record.id}
      field={field}
      value={`${value}`}
      subRecords={subRecords}
    />
  );
}

export const RecordBrowser: FC = () => {
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
                    <Link to={`/browse/records/${ownerRecord.id}`}>
                      Back to Parent Record
                    </Link>
                  </div>
                ) : (
                  <></>
                )}
                {form?.fields.map((f) =>
                  mapRecordField(record, f, subRecords?.byFieldId[f.id]),
                )}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

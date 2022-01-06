import React, { FC } from 'react';
import { FieldDefinition, FieldValue } from 'core-api-client';

import { FormValue } from '../../reducers/recorder';
import { FieldEditor } from '../../components/FieldEditor/FieldEditor';

type Props = {
  fields: FieldDefinition[];
  values: FieldValue[];
  onAddSubRecord: (ownerFieldId: string) => void;
  onChangeValue: (key: string, value: any) => void;
  onSaveRecord: () => void;
  onSelectSubRecord: (subRecordId: string) => void;
  subRecords: { [key: string]: FormValue[] };
};

export const RecordEditorComponent: FC<Props> = (props) => {
  const {
    fields,
    onAddSubRecord,
    onSelectSubRecord,
    onSaveRecord,
    subRecords,
    onChangeValue,
    values,
  } = props;
  if (!fields) {
    return <></>;
  }

  return (
    <div className="flex-grow-1 w-100 h-100 bg-dark text-light py-3 overflow-scroll">
      <div className="container-fluid">
        <div className="row justify-content-center">
          <div className="col-12 col-lg-8">
            <h4 className="mb-4">Add record</h4>
            <div className="card bg-dark text-light border-secondary">
              <div className="card-body">
                {fields.map((field) => {
                  const fieldValue = values.find((v) => v.fieldId === field.id);
                  const value = fieldValue?.value ? fieldValue.value : '';
                  const handleValueChange = (v: any) => {
                    onChangeValue(field.id, v);
                  };
                  const handleAddSubRecordWrapper = () => {
                    onAddSubRecord(field.id);
                  };
                  return (
                    <FieldEditor
                      key={field.id}
                      field={field}
                      value={value}
                      onChange={handleValueChange}
                      subRecords={subRecords[field.id]}
                      onSelectSubRecord={onSelectSubRecord}
                      onAddSubRecord={handleAddSubRecordWrapper}
                    />
                  );
                })}
                <div className="my-3">
                  <button
                    onClick={() => onSaveRecord()}
                    className="btn btn-primary"
                  >
                    Save Record
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

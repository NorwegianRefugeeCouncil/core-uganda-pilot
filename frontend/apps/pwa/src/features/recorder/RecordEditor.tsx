import React, { FC } from 'react';
import { FieldDefinition, FieldValue } from 'core-api-client';

import { FormValue } from '../../reducers/recorder';

import { FieldEditor } from './FieldEditor';

export type RecordEditorProps = {
  fields: FieldDefinition[];
  values: FieldValue[];
  setValue: (key: string, value: any) => void;
  selectSubRecord: (subRecordId: string) => void;
  addSubRecord: (ownerFieldId: string) => void;
  subRecords: { [key: string]: FormValue[] };
  saveRecord: () => void;
};

export const RecordEditor: FC<RecordEditorProps> = (props) => {
  const {
    fields,
    addSubRecord,
    selectSubRecord,
    saveRecord,
    subRecords,
    setValue,
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
                  const setValueWrapper = (v: any) => {
                    setValue(field.id, v);
                  };
                  const addSubRecordWrapper = () => {
                    addSubRecord(field.id);
                  };
                  return (
                    <FieldEditor
                      key={field.id}
                      field={field}
                      value={value}
                      setValue={setValueWrapper}
                      subRecords={subRecords[field.id]}
                      selectSubRecord={selectSubRecord}
                      addSubRecord={addSubRecordWrapper}
                    />
                  );
                })}
                <div className="my-3">
                  <button
                    onClick={() => saveRecord()}
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

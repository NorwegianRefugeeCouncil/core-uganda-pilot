import React, { FC } from 'react';
import { FieldDefinition, FieldKind, FieldValue } from 'core-api-client';
import { FieldErrors } from 'react-hook-form';
import { ErrorMessage } from '@hookform/error-message';

import { FieldEditor } from '../../components/FieldEditor/FieldEditor';
import { FormValue } from '../../reducers/Recorder/types';
import { FieldLabel } from '../../components/formFields/FieldLabel';
import { FieldDescription } from '../../components/formFields/FieldDescription';

type Props = {
  fields: FieldDefinition[];
  values: FieldValue[];
  onAddSubRecord: (ownerFieldId: string) => void;
  onChangeValue: (key: string, value: any) => void;
  onSaveRecord: () => void;
  onSelectSubRecord: (subRecordId: string) => void;
  subRecords: { [key: string]: FormValue[] };
  errors: FieldErrors;
};

export const RecordEditorComponent: FC<Props> = ({
  fields,
  onAddSubRecord,
  onSelectSubRecord,
  onSaveRecord,
  subRecords,
  onChangeValue,
  values,
  errors,
}) => {
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
              <div className="card-body" id="values">
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
                    <div className="form-group mb-2" key={field.id}>
                      {(field.fieldType !== FieldKind.Checkbox ||
                        field.fieldType !== FieldKind.SubForm) && (
                        <FieldLabel fieldDefinition={field} />
                      )}
                      <FieldEditor
                        field={field}
                        value={value}
                        onChange={handleValueChange}
                        subRecords={subRecords[field.id]}
                        onSelectSubRecord={onSelectSubRecord}
                        onAddSubRecord={handleAddSubRecordWrapper}
                        errors={errors}
                      />

                      <FieldDescription
                        text={field.description}
                        fieldId={field.id}
                      />
                      <div className="invalid-feedback" id="errorMessages">
                        <ErrorMessage
                          errors={errors}
                          name={`values.${field.id}`}
                        />
                      </div>
                    </div>
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

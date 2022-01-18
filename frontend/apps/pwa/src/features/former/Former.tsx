import React, { FC } from 'react';
import _ from 'lodash';

import { FormName } from './FormName';
import { FormerField } from './FormerField';
import { FormerProps } from './types';
import { FieldSections } from './FieldSections';

export const Former: FC<FormerProps> = (props) => {
  const {
    formName,
    setFormName,
    fields,
    selectedFieldId,
    saveForm,
    ownerFormName,
    setSelectedField,
    setFieldName,
    setFieldOption,
    addOption,
    removeOption,
    setFieldDescription,
    setFieldIsKey,
    setFieldRequired,
    openSubForm,
    saveField,
    cancelField,
    setFieldReferencedDatabaseId,
    setFieldReferencedFormId,
    error,
  } = props;

  const mappedErrors = error && _.keyBy(error, 'field');
  console.log('MAPPED', mappedErrors);

  const selectedField = selectedFieldId
    ? fields.find((f) => f.id === selectedFieldId)
    : undefined;

  function formHeader() {
    const name = formName || ownerFormName || '';
    return <FormName formName={name} setFormName={setFormName} />;
  }

  if (selectedField) {
    return (
      <div className="flex-grow-1 bg-dark text-light overflow-scroll">
        <div className="container-fluid mt-4">
          <div className="row">
            <div className="col-12 col-md-8 offset-md-1">
              <h3>Add Form</h3>
              <h6>
                {ownerFormName ? (
                  <div className="mb-2">
                    Editing child form of {ownerFormName}
                  </div>
                ) : (
                  <></>
                )}
              </h6>
            </div>
          </div>
          <div className="row mt-3">
            <div className="col-10 col-md-8 offset-md-1">
              {ownerFormName == null && formHeader()}
              <FormerField
                key={selectedField.id}
                isSelected={selectedField.id === selectedFieldId}
                selectField={() => setSelectedField(selectedField.id)}
                fieldType={selectedField.type}
                fieldOptions={selectedField.options}
                setFieldOption={(i: number, value: string) =>
                  setFieldOption(selectedField.id, i, value)
                }
                addOption={() => addOption(selectedField.id)}
                removeOption={(i: number) => removeOption(selectedField.id, i)}
                fieldName={selectedField.name}
                setFieldName={(name) => setFieldName(selectedField.id, name)}
                fieldRequired={selectedField.required}
                setFieldRequired={(req) =>
                  setFieldRequired(selectedField.id, req)
                }
                fieldIsKey={selectedField.key}
                setFieldIsKey={(isKey) =>
                  setFieldIsKey(selectedField.id, isKey)
                }
                fieldDescription={selectedField.description}
                setFieldDescription={(d) =>
                  setFieldDescription(selectedField.id, d)
                }
                openSubForm={() => openSubForm(selectedField.id)}
                cancel={() => cancelField(selectedField.id)}
                saveField={() => saveField(selectedField.id)}
                referencedDatabaseId={selectedField.referencedDatabaseId}
                referencedFormId={selectedField.referencedFormId}
                setReferencedDatabaseId={(d) =>
                  setFieldReferencedDatabaseId(selectedField.id, d)
                }
                setReferencedFormId={(d) =>
                  setFieldReferencedFormId(selectedField.id, d)
                }
              />
            </div>
            <div className="col-2">
              <button
                className="btn btn-primary w-100"
                onClick={() => saveForm()}
              >
                Save
              </button>
            </div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="h-100 w-100 bg-dark text-light flex-grow-1 overflow-scroll">
      <div className="container mt-4">
        <div className="row">
          <div className="col-8 offset-2">
            <h3>Add Form</h3>
            <h6>
              {ownerFormName ? (
                <div className="mb-2 p-2 border-secondary">
                  Editing child form of {ownerFormName}
                </div>
              ) : (
                <></>
              )}
            </h6>
          </div>
        </div>
        <div className="row mt-3">
          <div className="col-6 offset-2">
            {ownerFormName == null && formHeader()}
            <FieldSections formerProps={props} />
            {mappedErrors.fields &&
              error.map((e, i) => (
                <div
                  key={`${e.reason}_${i}`}
                  className="is-invalid invalid-feedback"
                >
                  {e?.message}
                </div>
              ))}
          </div>
          <div className="col-2">
            <button
              className="btn btn-primary w-100"
              onClick={() => saveForm()}
            >
              Save
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

import React, { FC } from 'react';

import { FormerField } from './FormerField';
import { FormerProps } from './types';
import { FieldSections } from './FieldSections';
import { FormName } from './FormName';

export const Former: FC<FormerProps> = (props) => {
  const {
    addField,
    addOption,
    cancelField,
    errors,
    fields,
    formName,
    openSubForm,
    ownerFormName,
    removeOption,
    saveField,
    saveForm,
    selectedFieldId,
    setFieldDescription,
    setFieldIsKey,
    setFieldName,
    setFieldOption,
    setFieldReferencedDatabaseId,
    setFieldReferencedFormId,
    setFieldRequired,
    setFormName,
    setSelectedField,
  } = props;

  const selectedField = selectedFieldId
    ? fields.find((f) => f.id === selectedFieldId)
    : undefined;

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
              {ownerFormName == null && (
                <FormName
                  setFormName={setFormName}
                  formName={formName}
                  errors={errors}
                />
              )}
              <FormerField
                key={selectedField.id}
                isSelected={selectedField.id === selectedFieldId}
                selectField={() => setSelectedField(selectedField.id)}
                fieldType={selectedField.fieldType}
                fieldOptions={selectedField.options}
                errors={selectedField.errors}
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
            {ownerFormName == null && (
              <FormName
                setFormName={setFormName}
                formName={formName}
                errors={errors}
              />
            )}
            <FieldSections
              addField={addField}
              addOption={addOption}
              cancelField={(fieldId: string) => cancelField(fieldId)}
              errors={errors}
              fields={fields}
              openSubForm={openSubForm}
              removeOption={removeOption}
              saveField={saveField}
              selectedFieldId={selectedFieldId}
              setFieldDescription={setFieldDescription}
              setFieldIsKey={setFieldIsKey}
              setFieldName={setFieldName}
              setFieldOption={setFieldOption}
              setFieldReferencedDatabaseId={setFieldReferencedDatabaseId}
              setFieldReferencedFormId={setFieldReferencedFormId}
              setFieldRequired={setFieldRequired}
              setSelectedField={setSelectedField}
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
};

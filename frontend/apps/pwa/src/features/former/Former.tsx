import React, { FC } from 'react';
import { ErrorMessage } from '@hookform/error-message';

import { FormTypeControl } from './FormTypeControl';
import { FormerProps } from './types';
import { FieldSections } from './FieldSections';
import { FormName } from './FormName';

export const Former: FC<FormerProps> = (props) => {
  const {
    formId,
    formType,
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
            <>
              {ownerFormName == null && (
                <FormName
                  errors={errors}
                  formName={formName}
                  setFormName={setFormName}
                />
              )}
              <FormTypeControl formId={formId} formType={formType} />
            </>
            <FieldSections
              addField={addField}
              addOption={addOption}
              cancelField={(fieldId: string) => cancelField(fieldId)}
              errors={errors}
              fields={fields}
              formType={formType}
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
            <div className="text-danger">
              <ErrorMessage errors={errors} name="fields" />
            </div>
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

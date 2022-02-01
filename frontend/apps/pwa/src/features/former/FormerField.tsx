import { ErrorMessage } from '@hookform/error-message';
import {
  FieldDefinition,
  FieldKind,
  FormType,
  SelectOption,
} from 'core-api-client';
import React, { FC } from 'react';
import { FieldErrors } from 'react-hook-form';

import DatabasePickerContainer from '../../components/DatabasePicker';
import FormPickerContainer from '../../components/FormPicker';
import { Form } from '../../reducers/Former/types';

import validation from './validation';

type FormerFieldProps = {
  formType: FormType;
  addOption: () => void;
  cancel: () => void;
  errors: FieldErrors<Form & { selectedField?: FieldDefinition }>;
  fieldDescription: string;
  fieldIsKey: boolean;
  fieldName: string;
  fieldOptions?: SelectOption[];
  fieldRequired: boolean;
  fieldType: FieldKind;
  isSelected: boolean;
  openSubForm: () => void;
  referencedDatabaseId: string | undefined;
  referencedFormId: string | undefined;
  register: any;
  removeOption: (index: number) => void;
  revalidate: any;
  saveField: () => void;
  selectField: () => void;
  setFieldDescription: (description: string) => void;
  setFieldIsKey: (isKey: boolean) => void;
  setFieldName: (fieldName: string) => void;
  setFieldOption: (i: number, value: string) => void;
  setFieldRequired: (required: boolean) => void;
  setReferencedDatabaseId: (databaseId: string) => void;
  setReferencedFormId: (formId: string) => void;
};

export const FormerField: FC<FormerFieldProps> = (props) => {
  const {
    formType,
    addOption,
    cancel,
    errors = undefined,
    fieldDescription,
    fieldIsKey,
    fieldName,
    fieldOptions = [],
    fieldRequired,
    fieldType,
    isSelected,
    openSubForm,
    referencedDatabaseId,
    referencedFormId,
    register,
    removeOption,
    revalidate,
    saveField,
    selectField,
    setFieldDescription,
    setFieldIsKey,
    setFieldName,
    setFieldOption,
    setFieldRequired,
    setReferencedDatabaseId,
    setReferencedFormId,
  } = props;

  React.useEffect(() => {
    if (fieldType === FieldKind.Checkbox) {
      setFieldRequired(true);
      setFieldIsKey(false);
    }
  }, [fieldType]);

  const requiredDisabled = fieldIsKey || fieldType === FieldKind.Checkbox;
  const isKeyDisabled = fieldType === FieldKind.Checkbox;

  const registerSelectedFieldName = register(
    'selectedField.name',
    validation.selectedField.name,
  );
  const registerSelectedFieldOptionsSingle = register(
    'selectedField.fieldType.singleSelect.options',
    validation.selectedField.fieldType.singleSelect.options,
  );
  const registerSelectedFieldOptionsMulti = register(
    'selectedField.fieldType.multiSelect.options',
    validation.selectedField.fieldType.multiSelect.options,
  );

  if (!isSelected) {
    return (
      <div>
        <div
          onClick={() => selectField()}
          style={{ cursor: 'pointer' }}
          className={`card bg-dark text-light border-light mb-2 ${
            errors?.selectedField?.name ? 'is-invalid' : ''
          }`}
          id="name"
          aria-describedby="nameFeedback"
        >
          <div className="card-body p-3">
            <div className="d-flex flex-row">
              <span className="flex-grow-1">{fieldName}</span>
              <small className="text-uppercase">{fieldType}</small>
            </div>
          </div>
        </div>
        <ErrorMessage errors={errors} name="selectedField" />
      </div>
    );
  }

  return (
    <div className="card text-dark mb-2">
      <div className="card-body">
        {/* Form Title */}

        <h6 className="card-title text-uppercase">{fieldType}</h6>

        <div className="row">
          {/* Left Hand Side Section */}

          <div className="col-8">
            {/* Form Name */}

            <div className="form-group mb-2">
              <label className="form-label" htmlFor="name">
                Field Name
              </label>
              <input
                className={`form-control ${
                  errors?.selectedField?.name ? 'is-invalid' : ''
                }`}
                type="text"
                id="name"
                value={fieldName || ''}
                {...registerSelectedFieldName}
                onChange={(event) => {
                  setFieldName(event.target.value);
                  return registerSelectedFieldName.onChange(event);
                }}
              />
              <div className="invalid-feedback" id="nameFeedback">
                <ErrorMessage errors={errors} name="selectedField.name" />
              </div>
            </div>

            {/* Options */}

            {(fieldType === FieldKind.SingleSelect ||
              fieldType === FieldKind.MultiSelect) && (
              <div className="form-group mb-2">
                <div
                  className={`d-flex justify-content-between align-items-center mb-2 ${
                    errors?.selectedField?.fieldType?.singleSelect?.options
                      ? 'is-invalid'
                      : ''
                  }`}
                  id="options"
                  aria-describedby="optionsFeedback"
                >
                  <label className="form-label" htmlFor="fieldName">
                    Field Options
                  </label>
                  <button
                    type="button"
                    className="btn btn-outline-primary"
                    onClick={() => addOption()}
                  >
                    Add option
                  </button>
                </div>
                {fieldOptions?.map((opt, i) => (
                  <div key={i} className="d-flex mb-2">
                    <input
                      className="form-control me-3"
                      id={`fieldOption-${i}`}
                      type="text"
                      value={opt ? opt.name : ''}
                      {...registerSelectedFieldOptionsSingle}
                      onChange={(event) => {
                        setFieldOption(i, event.target.value);
                        registerSelectedFieldOptionsSingle.onChange(event);
                        registerSelectedFieldOptionsMulti.onChange(event);
                      }}
                    />
                    <button
                      type="button"
                      className="btn btn-outline-danger"
                      onClick={() => removeOption(i)}
                    >
                      <i className="bi bi-x" />
                    </button>
                  </div>
                ))}

                <div className="invalid-feedback" id="nameFeedback">
                  {fieldType === FieldKind.SingleSelect && (
                    <ErrorMessage
                      errors={errors}
                      name="selectedField.fieldType.singleSelect.options"
                    />
                  )}
                  {fieldType === FieldKind.MultiSelect && (
                    <ErrorMessage
                      errors={errors}
                      name="selectedField.fieldType.multiSelect.options"
                    />
                  )}
                </div>
              </div>
            )}

            {/* Form Description */}

            <div className="form-group mb-2">
              <label className="form-label" htmlFor="fieldName">
                Description
              </label>
              <textarea
                className="form-control"
                id="fieldDescription"
                onChange={(event) => setFieldDescription(event.target.value)}
                value={fieldDescription || ''}
              />
            </div>

            {/* Open Subform Button */}

            {fieldType === FieldKind.SubForm ? (
              <button
                className="btn btn-primary"
                onClick={() => openSubForm()}
                id="subformFields"
                aria-describedby="subformFieldsFeedback"
              >
                Open Sub Form
              </button>
            ) : (
              <></>
            )}

            {/* Configure Reference Field */}

            {fieldType === FieldKind.Reference ? (
              <div>
                <div className="form-group mb-2">
                  <label className="form-label">Database</label>
                  <DatabasePickerContainer
                    setDatabaseId={setReferencedDatabaseId}
                    databaseId={referencedDatabaseId}
                  />
                </div>
                <div className="form-group">
                  <label className="form-label">Form</label>
                  <FormPickerContainer
                    databaseId={referencedDatabaseId}
                    formId={referencedFormId}
                    setFormId={setReferencedFormId}
                    isRecipientKey={
                      fieldIsKey && formType === FormType.RecipientFormType
                    }
                  />
                </div>
              </div>
            ) : (
              <></>
            )}
          </div>

          {/* Right Hand Side Section */}

          <div className="col-4">
            {/* Required Checkbox */}

            <div className="form-check">
              <input
                disabled={requiredDisabled}
                className="form-check-input"
                type="checkbox"
                value=""
                onChange={() => setFieldRequired(!fieldRequired)}
                checked={fieldRequired}
                id="required"
              />
              <label className="form-check-label" htmlFor="required">
                Required
              </label>
            </div>

            {/* Key Checkbox */}

            <div className="form-check">
              <input
                className="form-check-input"
                type="checkbox"
                value=""
                onChange={() => setFieldIsKey(!fieldIsKey)}
                checked={fieldIsKey}
                disabled={isKeyDisabled}
                id="key"
              />
              <label className="form-check-label" htmlFor="key">
                Key
              </label>
            </div>
          </div>
        </div>
      </div>
      <ErrorMessage errors={errors} name="selectedField" />
      <div className="card-footer">
        <button
          onClick={async () => {
            const valid = await revalidate('selectedField');
            if (valid) {
              await saveField();
            }
          }}
          className="btn btn-primary me-2 shadow"
        >
          Save
        </button>
        <button onClick={() => cancel()} className="btn btn-secondary shadow">
          Delete
        </button>
      </div>
    </div>
  );
};

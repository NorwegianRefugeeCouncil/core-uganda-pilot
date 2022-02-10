import { ErrorMessage } from '@hookform/error-message';
import { FieldKind, FormType } from 'core-api-client';
import React, { FC } from 'react';
import { FieldErrors, useFormContext } from 'react-hook-form';
import { UseFormTrigger } from 'react-hook-form/dist/types/form';

import DatabasePickerContainer from '../../components/DatabasePicker';
import FormPickerContainer from '../../components/FormPicker';
import { FormField, ValidationForm } from '../../reducers/Former/types';

import { registeredValidation } from './validation';

type FormerFieldProps = {
  addOption: () => void;
  cancel: () => void;
  errors: FieldErrors<ValidationForm>;
  field: FormField;
  formType: FormType;
  isSelected: boolean;
  openSubForm: () => void;
  removeOption: (index: number) => void;
  revalidate: UseFormTrigger<ValidationForm>;
  saveField: (field: FormField) => void;
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
    addOption,
    cancel,
    errors = undefined,
    field,
    formType,
    isSelected,
    openSubForm,
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

  const {
    description,
    fieldType,
    key,
    name,
    options,
    referencedDatabaseId,
    referencedFormId,
    required,
  } = field;

  const { register } = useFormContext();

  React.useEffect(() => {
    if (fieldType === FieldKind.Checkbox) {
      setFieldRequired(true);
      setFieldIsKey(false);
    }
  }, [fieldType]);

  const requiredDisabled = key || fieldType === FieldKind.Checkbox;
  const isKeyDisabled = fieldType === FieldKind.Checkbox;

  const registerSelectedFieldName = register(
    'selectedField.name',
    registeredValidation.selectedField.name,
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
          aria-describedby="errorMessages"
        >
          <div className="card-body p-3">
            <div className="d-flex flex-row">
              <span className="flex-grow-1">{name}</span>
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
                value={name || ''}
                {...registerSelectedFieldName}
                onChange={(event) => {
                  setFieldName(event.target.value);
                  return registerSelectedFieldName.onChange(event);
                }}
              />
              <div className="invalid-feedback" id="errorMessages">
                <ErrorMessage errors={errors} name="selectedField.name" />
              </div>
            </div>

            {/* Options */}

            {(fieldType === FieldKind.SingleSelect ||
              fieldType === FieldKind.MultiSelect) && (
              <div className="form-group mb-2">
                <div
                  className={`d-flex justify-content-between align-items-center mb-2 ${
                    errors?.selectedField?.fieldType?.singleSelect?.options ||
                    errors?.selectedField?.fieldType?.multiSelect?.options
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
                {options?.map((opt, i) => {
                  const registerOption = register(
                    `selectedField.fieldType.${fieldType}.options.${i}`,
                    registeredValidation.selectedField.fieldType.select.option
                      .name,
                  );

                  const optionErrors =
                    fieldType === FieldKind.SingleSelect
                      ? errors?.selectedField?.fieldType?.singleSelect?.options
                      : errors?.selectedField?.fieldType?.multiSelect?.options;

                  const isInvalid = optionErrors && !!optionErrors[i];
                  return (
                    <div key={i} className="form-group mb-2 d-inline-flex">
                      <div className="input-group mb-2">
                        <input
                          className={`form-control me-3 ${
                            isInvalid ? 'is-invalid' : ''
                          }`}
                          id={`fieldOption-${i}`}
                          type="text"
                          value={opt ? opt.name : ''}
                          {...registerOption}
                          onChange={(event) => {
                            setFieldOption(i, event.target.value);
                            return registerOption.onChange(event);
                          }}
                          aria-describedby="errorMessages"
                        />
                        <button
                          type="button"
                          className="btn btn-outline-danger"
                          onClick={() => removeOption(i)}
                        >
                          <i className="bi bi-x" />
                        </button>
                        <div className="invalid-feedback" id="errorMessages">
                          <ErrorMessage
                            errors={errors}
                            name={`selectedField.fieldType.${fieldType}.options.${i}`}
                          />
                        </div>
                      </div>
                    </div>
                  );
                })}

                <div className="invalid-feedback" id="errorMessages">
                  <ErrorMessage
                    errors={errors}
                    name={`selectedField.fieldType.${fieldType}.options`}
                  />
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
                value={description || ''}
              />
            </div>

            {/* Open Subform Button */}

            {fieldType === FieldKind.SubForm ? (
              <button
                className="btn btn-primary"
                onClick={async () => {
                  const valid = await revalidate('selectedField');
                  if (valid) openSubForm();
                }}
                id="subformFields"
                aria-describedby="subformFieldsFeedback"
              >
                Open Sub Form
              </button>
            ) : (
              <></>
            )}

            <div className="text-danger">
              <ErrorMessage
                errors={errors}
                name="selectedField.fieldType.subForm"
              />
            </div>

            {/* Configure Reference Field */}

            {fieldType === FieldKind.Reference ? (
              <div>
                <div className="form-group mb-2">
                  <label className="form-label">Database</label>
                  <DatabasePickerContainer
                    setDatabaseId={setReferencedDatabaseId}
                    databaseId={referencedDatabaseId}
                    errors={errors}
                  />
                </div>
                <div className="form-group">
                  <label className="form-label">Form</label>
                  <FormPickerContainer
                    databaseId={referencedDatabaseId}
                    formId={referencedFormId}
                    setFormId={setReferencedFormId}
                    isRecipientKey={
                      key && formType === FormType.RecipientFormType
                    }
                    errors={errors}
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
                onChange={() => setFieldRequired(!required)}
                checked={required}
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
                onChange={() => setFieldIsKey(!key)}
                checked={key}
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

      <div className="invalid-feedback" id="errorMessages">
        <ErrorMessage errors={errors} name="selectedField" />
      </div>

      <div className="card-footer">
        <button
          onClick={async () => {
            await saveField(field);
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

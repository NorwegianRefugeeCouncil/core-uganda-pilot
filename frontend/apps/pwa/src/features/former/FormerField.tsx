import { FieldKind, FormType, SelectOption } from 'core-api-client';
import React, { FC } from 'react';

import DatabasePickerContainer from '../../components/DatabasePicker';
import FormPickerContainer from '../../components/FormPicker';
import { ErrorMessage } from '../../types/errors';

type FormerFieldProps = {
  formType: FormType;
  addOption: () => void;
  cancel: () => void;
  errors: ErrorMessage | undefined;
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
  removeOption: (index: number) => void;
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
    fieldOptions,
    fieldRequired,
    fieldType,
    isSelected,
    openSubForm,
    referencedDatabaseId,
    referencedFormId,
    removeOption,
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
    if (fieldType === FieldKind.Boolean) {
      setFieldRequired(true);
      setFieldIsKey(false);
    }
  }, [fieldType]);

  const requiredDisabled = fieldIsKey || fieldType === FieldKind.Boolean;
  const isKeyDisabled = fieldType === FieldKind.Boolean;

  if (!isSelected) {
    return (
      <div>
        <div
          onClick={() => selectField()}
          style={{ cursor: 'pointer' }}
          className={`card bg-dark text-light border-light mb-2 ${
            errors?.name ? 'is-invalid' : ''
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
        {errors?.name && (
          <div className="invalid-feedback is-invalid" id="nameFeedback">
            {errors?.name}
          </div>
        )}
        {errors?.fieldType?.singleSelect?.options && (
          <div className="invalid-feedback is-invalid" id="optionsFeedback">
            {errors?.fieldType?.singleSelect?.options}
          </div>
        )}
        {errors?.fieldType?.multiSelect?.options && (
          <div className="invalid-feedback is-invalid" id="optionsFeedback">
            {errors?.fieldType?.multiSelect?.options}
          </div>
        )}
        {errors?.fieldType?.subForm?.fields && (
          <div
            className="invalid-feedback is-invalid"
            id="subformFieldsFeedback"
          >
            {errors?.fieldType?.subForm?.fields}
          </div>
        )}
      </div>
    );
  }

  return (
    <div className="card text-dark">
      <div className="card-body">
        {/* Form Title */}

        <h6 className="card-title text-uppercase">{fieldType}</h6>

        <div className="row">
          {/* Left Hand Side Section */}

          <div className="col-8">
            {/* Form Name */}

            <div className="form-group mb-2">
              <label className="form-label" htmlFor="fieldName">
                Field Name
              </label>
              <input
                className="form-control"
                type="text"
                value={fieldName || ''}
                onChange={(event) => setFieldName(event.target.value)}
              />
            </div>

            {/* Options */}

            {(fieldType === FieldKind.SingleSelect ||
              fieldType === FieldKind.MultiSelect) && (
              <div className="form-group mb-2">
                <div
                  className={`d-flex justify-content-between align-items-center mb-2 ${
                    errors?.fieldType?.singleSelect?.options ? 'is-invalid' : ''
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
                      onChange={(event) =>
                        setFieldOption(i, event.target.value)
                      }
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
      <div className="card-footer">
        <button
          onClick={() => saveField()}
          className="btn btn-primary me-2 shadow"
        >
          Save
        </button>
        <button onClick={() => cancel()} className="btn btn-secondary shadow">
          Cancel
        </button>
      </div>
    </div>
  );
};

import React, { FC, useCallback, useState } from 'react';
import { FieldKind } from 'core-api-client';

import { FormField } from '../../reducers/former';

import { FormerField } from './Field';
import { FormName } from './FormName';
import { FieldTypePicker } from './FieldTypePicker';

type FormerProps = {
  formName: string;
  setFormName: (formName: string) => void;
  fields: FormField[];
  fieldOptions?: string[];
  setFieldOption: (fieldId: string, i: number, value: string) => void;
  addOption: (fieldId: string) => void;
  removeOption: (fieldId: string, index: number) => void;
  selectedFieldId: string | undefined;
  setSelectedField: (fieldId: string | undefined) => void;
  addField: (kind: FieldKind) => void;
  setFieldRequired: (fieldId: string, required: boolean) => void;
  setFieldIsKey: (fieldId: string, isKey: boolean) => void;
  setFieldName: (fieldId: string, name: string) => void;
  setFieldDescription: (fieldId: string, description: string) => void;
  setFieldReferencedDatabaseId: (fieldId: string, databaseId: string) => void;
  setFieldReferencedFormId: (fieldId: string, formId: string) => void;
  openSubForm: (fieldId: string) => void;
  saveField: (fieldId: string) => void;
  saveForm: () => void;
  ownerFormName: string | undefined;
  cancelField: (fieldId: string) => void;
};

function mapField(f: FormField, props: FormerProps) {
  const {
    selectedFieldId,
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
  } = props;

  return (
    <FormerField
      key={f.id}
      isSelected={f.id === selectedFieldId}
      selectField={() => setSelectedField(f.id)}
      fieldType={f.type}
      fieldOptions={f.options}
      setFieldOption={(i: number, value: string) =>
        setFieldOption(f.id, i, value)
      }
      addOption={() => addOption(f.id)}
      removeOption={(i: number) => removeOption(f.id, i)}
      fieldName={f.name}
      setFieldName={(name) => setFieldName(f.id, name)}
      fieldRequired={f.required}
      setFieldRequired={(req) => setFieldRequired(f.id, req)}
      fieldIsKey={f.key}
      setFieldIsKey={(isKey) => setFieldIsKey(f.id, isKey)}
      fieldDescription={f.description}
      setFieldDescription={(d) => setFieldDescription(f.id, d)}
      openSubForm={() => openSubForm(f.id)}
      cancel={() => cancelField(f.id)}
      saveField={() => saveField(f.id)}
      referencedDatabaseId={f.referencedDatabaseId}
      referencedFormId={f.referencedFormId}
      setReferencedDatabaseId={(d) => setFieldReferencedDatabaseId(f.id, d)}
      setReferencedFormId={(d) => setFieldReferencedFormId(f.id, d)}
    />
  );
}

export const Former: FC<FormerProps> = (props) => {
  const {
    formName,
    setFormName,
    fields,
    selectedFieldId,
    addField,
    saveForm,
    ownerFormName,
  } = props;

  const [isAddingField, setIsAddingField] = useState(false);

  const selectedField = selectedFieldId
    ? fields.find((f) => f.id === selectedFieldId)
    : undefined;

  function formHeader() {
    const name = formName || ownerFormName || '';
    return <FormName formName={name} setFormName={setFormName} />;
  }

  function addFieldButton() {
    return (
      <div>
        <button
          className="btn btn-primary my-2 mb-3 w-100"
          onClick={() => setIsAddingField(true)}
        >
          Add field
        </button>
      </div>
    );
  }

  const fieldSections = useCallback(() => {
    if (isAddingField) {
      return (
        <FieldTypePicker
          onCancel={() => {
            setIsAddingField(false);
          }}
          onSubmit={(fieldKind) => {
            addField(fieldKind);
            setIsAddingField(false);
          }}
        />
      );
    }
    return (
      <div>
        {addFieldButton()}
        {fields.map((f) => mapField(f, props))}
      </div>
    );
  }, [addField, fields, isAddingField, props]);

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
              {mapField(selectedField, props)}
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
            {fieldSections()}
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

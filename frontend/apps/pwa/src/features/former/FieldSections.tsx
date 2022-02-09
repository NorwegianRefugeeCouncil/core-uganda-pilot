import React, { useState } from 'react';

import { FormField } from '../../reducers/Former/types';

import { FieldTypePicker } from './FieldTypePicker';
import { FormerField } from './FormerField';
import { FormerProps } from './types';

type Props = Omit<
  FormerProps,
  | 'fieldOptions'
  | 'formName'
  | 'ownerFormName'
  | 'saveForm'
  | 'setFormName'
  | 'formId'
  | 'invalid'
>;

export const FieldSections: React.FC<Props> = ({
  formType,
  addField,
  addOption,
  cancelField,
  errors,
  fields,
  openSubForm,
  register,
  removeOption,
  revalidate,
  saveField,
  selectedFieldId,
  setFieldDescription,
  setFieldIsKey,
  setFieldName,
  setFieldOption,
  setFieldReferencedDatabaseId,
  setFieldReferencedFormId,
  setFieldRequired,
  setSelectedField,
}) => {
  const [isAddingField, setIsAddingField] = useState(false);

  const handleAddOption = (f: FormField) => () => addOption(f.id);
  const handleCancel = (f: FormField) => () => cancelField(f.id);
  const handleOpenSubForm = (f: FormField) => () => openSubForm(f.id);
  const handleSelectField = (f: FormField) => () => setSelectedField(f.id);
  const handleSetFieldRequired = (f: FormField) => (required: boolean) =>
    setFieldRequired(f.id, required);
  const handleSetDataBaseId = (f: FormField) => (dataBaseId: string) =>
    setFieldReferencedDatabaseId(f.id, dataBaseId);
  const handleSetFormId = (f: FormField) => (formId: string) =>
    setFieldReferencedFormId(f.id, formId);
  const handleSetFieldOption = (f: FormField) => (i: number, value: string) =>
    setFieldOption(f.id, i, value);
  const handleSetFieldIsKey = (f: FormField) => (isKey: boolean) =>
    setFieldIsKey(f.id, isKey);
  const handleSetFieldName = (f: FormField) => (name: string) =>
    setFieldName(f.id, name);
  const handleSetFieldDescription = (f: FormField) => (d: string) =>
    setFieldDescription(f.id, d);
  const handleRemoveOption = (f: FormField) => (i: number) =>
    removeOption(f.id, i);

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
      {!selectedFieldId && (
        <div>
          <button
            className="btn btn-primary my-2 mb-3 w-100"
            onClick={() => {
              setIsAddingField(true);
            }}
          >
            Add field
          </button>
        </div>
      )}
      {fields.map((f: FormField) => (
        <FormerField
          addOption={handleAddOption(f)}
          cancel={handleCancel(f)}
          errors={errors}
          field={f}
          formType={formType}
          isSelected={f.id === selectedFieldId}
          key={f.id}
          openSubForm={handleOpenSubForm(f)}
          register={register}
          removeOption={handleRemoveOption(f)}
          revalidate={revalidate}
          saveField={saveField}
          selectField={handleSelectField(f)}
          setFieldDescription={handleSetFieldDescription(f)}
          setFieldIsKey={handleSetFieldIsKey(f)}
          setFieldName={handleSetFieldName(f)}
          setFieldOption={handleSetFieldOption(f)}
          setFieldRequired={handleSetFieldRequired(f)}
          setReferencedDatabaseId={handleSetDataBaseId(f)}
          setReferencedFormId={handleSetFormId(f)}
        />
      ))}
    </div>
  );
};

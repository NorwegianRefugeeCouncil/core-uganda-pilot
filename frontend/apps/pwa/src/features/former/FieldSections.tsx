import React, { useState } from 'react';

import { FormField } from '../../reducers/Former/types';

import { FieldTypePicker } from './FieldTypePicker';
import { FormerField } from './FormerField';
import { FormerProps } from './types';

type FieldSectionsProps = Omit<
  FormerProps,
  | 'fieldOptions'
  | 'formName'
  | 'ownerFormName'
  | 'saveForm'
  | 'setFormName'
  | 'formId'
>;

export const FieldSections: React.FC<FieldSectionsProps> = (props) => {
  const [isAddingField, setIsAddingField] = useState(false);

  const {
    formType,
    addField,
    addOption,
    cancelField,
    fields,
    openSubForm,
    removeOption,
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
  } = props;

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
      <div>
        <button
          className="btn btn-primary my-2 mb-3 w-100"
          onClick={() => setIsAddingField(true)}
        >
          Add field
        </button>
      </div>
      {fields.map((f: FormField) => (
        <FormerField
          key={f.id}
          formType={formType}
          addOption={() => addOption(f.id)}
          cancel={() => cancelField(f.id)}
          errors={f.errors}
          fieldDescription={f.description}
          fieldIsKey={f.key}
          fieldName={f.name}
          fieldOptions={f.options}
          fieldRequired={f.required}
          fieldType={f.fieldType}
          isSelected={f.id === selectedFieldId}
          openSubForm={() => openSubForm(f.id)}
          referencedDatabaseId={f.referencedDatabaseId}
          referencedFormId={f.referencedFormId}
          removeOption={(i: number) => removeOption(f.id, i)}
          saveField={() => saveField(f.id)}
          selectField={() => setSelectedField(f.id)}
          setFieldDescription={(d) => setFieldDescription(f.id, d)}
          setFieldIsKey={(isKey) => setFieldIsKey(f.id, isKey)}
          setFieldName={(name) => setFieldName(f.id, name)}
          setFieldOption={(i: number, value: string) =>
            setFieldOption(f.id, i, value)
          }
          setFieldRequired={(req) => setFieldRequired(f.id, req)}
          setReferencedDatabaseId={(d) => setFieldReferencedDatabaseId(f.id, d)}
          setReferencedFormId={(d) => setFieldReferencedFormId(f.id, d)}
        />
      ))}
    </div>
  );
};

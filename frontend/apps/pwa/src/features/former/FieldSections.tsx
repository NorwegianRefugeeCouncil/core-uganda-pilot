import React, { useState } from 'react';

import { FormField } from '../../reducers/former';

import { FieldTypePicker } from './FieldTypePicker';
import { FormerField } from './FormerField';
import { FormerProps } from './types';

type FieldSectionsProps = {
  formerProps: FormerProps;
};

export const FieldSections: React.FC<FieldSectionsProps> = ({
  formerProps,
}) => {
  const [isAddingField, setIsAddingField] = useState(false);

  const {
    fields,
    selectedFieldId,
    addField,
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
  } = formerProps;

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
      ))}
    </div>
  );
};

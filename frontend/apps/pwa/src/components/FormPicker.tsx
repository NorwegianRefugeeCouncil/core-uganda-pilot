import { FC, useCallback } from 'react';
import { FormDefinition, FormType } from 'core-api-client';

import { useDatabase, useForms } from '../app/hooks';

export type FormPickerProps = {
  forms: FormDefinition[];
  formId: string | undefined;
  disabled?: boolean;
  setFormId: (formId: string) => void;
};

export const FormPicker: FC<FormPickerProps> = (props) => {
  const { forms, formId, setFormId, disabled } = props;
  const hasForms = forms.length > 0;
  return (
    <div>
      <select
        disabled={disabled || !hasForms}
        onChange={(e) => setFormId(e.target.value)}
        value={formId || ''}
        className="form-select"
        aria-label="Select Form"
      >
        <option disabled value="">
          {hasForms ? 'Select Form' : 'No Forms'}
        </option>
        {forms.map((f) => {
          return <option value={f.id}>{f.name}</option>;
        })}
      </select>
    </div>
  );
};

export type FormPickerContainerProps = {
  databaseId: string | undefined;
  formId: string | undefined;
  isRecipientKey: boolean;
  setForm?: (form: FormDefinition | undefined) => void;
  setFormId?: (formId: string) => void;
};

const FormPickerContainer: FC<FormPickerContainerProps> = ({
  databaseId,
  formId,
  setFormId,
  setForm,
  isRecipientKey,
}) => {
  const database = useDatabase(databaseId);
  const allForms = useForms({ databaseId });
  const forms = isRecipientKey
    ? allForms.filter((f) => f.formType === FormType.RecipientFormType)
    : allForms;

  const callback = useCallback(
    (formId: string) => {
      if (setFormId) {
        setFormId(formId);
      }
      if (setForm) {
        const form = forms.find((f) => f.id === formId);
        setForm(form);
      }
    },
    [forms, setForm, setFormId],
  );

  if (!database) {
    return (
      <FormPicker
        disabled
        setFormId={(formId1) => {}}
        forms={[]}
        formId={formId}
      />
    );
  }

  return <FormPicker setFormId={callback} forms={forms} formId={formId} />;
};

export default FormPickerContainer;

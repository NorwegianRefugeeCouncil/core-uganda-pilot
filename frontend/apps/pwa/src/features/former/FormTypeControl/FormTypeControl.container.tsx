import * as React from 'react';
import { FormType } from 'core-api-client';

import { useAppDispatch } from '../../../app/hooks';
import former from '../../../reducers/Former';

import { FormTypeControlComponent } from './FormTypeControl.component';

type Props = {
  formId: string;
  formType: FormType;
};

export const FormTypeControlContainer: React.FC<Props> = ({
  formId,
  formType,
}) => {
  const dispatch = useAppDispatch();

  const handleFormTypeChange = (ft: FormType) => {
    dispatch(former.actions.setFormType({ formId, formType: ft }));
  };

  return (
    <FormTypeControlComponent
      formType={formType}
      onFormTypeChange={handleFormTypeChange}
    />
  );
};

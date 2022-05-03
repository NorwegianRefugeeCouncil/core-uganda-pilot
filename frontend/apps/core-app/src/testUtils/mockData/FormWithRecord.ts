import { FormType, FormWithRecord } from 'core-api-client';
import { Recipient } from 'core-api-client/src/types/client/Recipient';

import { makeField } from './FormField';
import { makeForm } from './Form';
import { makeRecord } from './Record';

export const makeFormWithRecord = (i: number): FormWithRecord<Recipient> => {
  const field = makeField(i, false, false, { text: {} });
  const form = makeForm(i, FormType.DefaultFormType, [field]);
  const record = makeRecord(i, form);

  return {
    form,
    record,
  };
};

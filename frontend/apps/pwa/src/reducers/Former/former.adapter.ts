import { createEntityAdapter } from '@reduxjs/toolkit';

import { Form } from './types';

export const adapter = createEntityAdapter<Form>({
  selectId: (form) => form.formId,
  sortComparer: (a, b) => a.name.localeCompare(b.name),
});

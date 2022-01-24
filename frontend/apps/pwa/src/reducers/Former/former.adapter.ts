import { createEntityAdapter, createSlice } from '@reduxjs/toolkit';

import { Form } from './types';

export const adapter = createEntityAdapter<Form>({
  // Assume IDs are stored in a field other than `book.id`
  selectId: (form) => form.formId,
  // Keep the "all IDs" array sorted based on book titles
  sortComparer: (a, b) => a.formId.localeCompare(b.formId),
});

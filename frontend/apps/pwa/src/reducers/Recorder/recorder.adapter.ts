import { createEntityAdapter } from '@reduxjs/toolkit';

import { FormValue } from './types';

export const adapter = createEntityAdapter<FormValue>({
  // Assume IDs are stored in a field other than `book.id`
  selectId: (record) => record.id,
  // Keep the "all IDs" array sorted based on book titles
  sortComparer: (a, b) => a.id.localeCompare(b.id),
});

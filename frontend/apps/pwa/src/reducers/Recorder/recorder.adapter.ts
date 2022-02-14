import { createEntityAdapter } from '@reduxjs/toolkit';

import { FormValue } from './types';

export const adapter = createEntityAdapter<FormValue>({
  selectId: (record) => record.id,
  sortComparer: (a, b) => a.id.localeCompare(b.id),
});

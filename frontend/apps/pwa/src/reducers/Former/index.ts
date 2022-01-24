import { createSlice } from '@reduxjs/toolkit';

import { adapter } from './former.adapter';
import { extraReducers, reducers } from './former.reducers';

export default createSlice({
  name: 'recorder',
  initialState: {
    ...adapter.getInitialState(),
    selectedFormId: '',
    selectedFieldId: undefined as string | undefined,
    selectedDatabaseId: undefined as string | undefined,
    selectedFolderId: undefined as string | undefined,
    savePending: false,
    saveSuccess: false,
    saveId: undefined as string | undefined,
    saveError: undefined,
  },
  reducers,
  extraReducers,
});

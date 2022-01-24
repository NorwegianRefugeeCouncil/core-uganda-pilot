import { createSlice } from '@reduxjs/toolkit';

import { reducers, extraReducers } from './recorder.reducers';
import { adapter } from './recorder.adapter';

export default createSlice({
  name: 'recorder',
  initialState: {
    ...adapter.getInitialState(),
    selectedRecordId: '',
    baseFormId: '',
    editingValues: {} as { [recordId: string]: { [key: string]: any } },
    saveError: undefined,
  },
  reducers,
  extraReducers,
});

import { Action, configureStore, ThunkAction } from '@reduxjs/toolkit';

import formReducer from '../reducers/form';
import folderReducer from '../reducers/folder';
import databaseReducer from '../reducers/database';
import recordsReducer from '../reducers/records';
import recorderSlice from '../reducers/Recorder';
import formerSlice from '../reducers/Former';

export const store = configureStore({
  reducer: {
    forms: formReducer,
    folders: folderReducer,
    databases: databaseReducer,
    records: recordsReducer,
    recorder: recorderSlice.reducer,
    former: formerSlice.reducer,
  },
  devTools: true,
});

export type AppDispatch = typeof store.dispatch;
export type RootState = ReturnType<typeof store.getState>;
export type AppThunk<ReturnType = void> = ThunkAction<
  ReturnType,
  RootState,
  unknown,
  Action<string>
>;

import {Action, configureStore, ThunkAction} from '@reduxjs/toolkit';
import formReducer from "../reducers/form";
import folderReducer from "../reducers/folder";
import databaseReducer from "../reducers/database";
import recordsReducer from "../reducers/records";
import recorerReducer from "../features/recorder/recorder.slice";
import formerReducer from "../features/former/former.slice";
import organizationsReducer from "../reducers/organizations";
import identityProvidersReducer from "../reducers/identityproviders";

export const store = configureStore({
    reducer: {
        forms: formReducer,
        folders: folderReducer,
        databases: databaseReducer,
        records: recordsReducer,
        recorder: recorerReducer,
        former: formerReducer,
        organizations: organizationsReducer,
        identityProviders: identityProvidersReducer,
    },
    devTools: true
});

export type AppDispatch = typeof store.dispatch;
export type RootState = ReturnType<typeof store.getState>;
export type AppThunk<ReturnType = void> = ThunkAction<ReturnType,
    RootState,
    unknown,
    Action<string>>;

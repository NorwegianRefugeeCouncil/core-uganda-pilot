import { combineReducers, configureStore } from '@reduxjs/toolkit';
import formDefinitionsReducer, { State } from '../reducers/formdefinitions';
import * as formBuilder from '@core/formbuilder';


export type RootState = {
  formDefinitions: State
}

const rootReducer = combineReducers({
  formDefinitions: formDefinitionsReducer,
  formBuilder: formBuilder.reducer
});

export const store = configureStore({
  reducer: rootReducer,
  devTools: true
});


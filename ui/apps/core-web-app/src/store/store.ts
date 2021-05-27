import { combineReducers, configureStore } from '@reduxjs/toolkit';
import formDefinitionsReducer, { State } from '../reducers/formdefinitions';


export type RootState = {
  formDefinitions: State
}

const rootReducer = combineReducers({
  formDefinitions: formDefinitionsReducer
})

export const store = configureStore({
  reducer: rootReducer,
  devTools: true
});


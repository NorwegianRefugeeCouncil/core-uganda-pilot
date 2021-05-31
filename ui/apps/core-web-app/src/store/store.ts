import { combineReducers, configureStore } from '@reduxjs/toolkit';
import formDefinitionsReducer, { State } from '../reducers/formdefinitions';
import * as formBuilder from '@core/formbuilder';
import { entitiesReducer } from '@core/api-client';
import { Observable } from 'rxjs';


export type RootState = {
  formDefinitions: State
}

const rootReducer = combineReducers({
  formDefinitions: formDefinitionsReducer,
  formBuilder: formBuilder.reducer,
  entities: entitiesReducer
});

export const store = configureStore({
  reducer: rootReducer,
  devTools: true
});

export const state$ = new Observable(subscriber => {
  const unsubscribe = store.subscribe(() => {
    subscriber.next(store.getState());
  });
  return () => unsubscribe();
});

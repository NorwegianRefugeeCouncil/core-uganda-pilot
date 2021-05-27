import { FieldType, FormElement } from '@core/api-client';
import { createAction, createReducer } from '@reduxjs/toolkit';

export const KEY = 'formBuilder';

export type State = {
  root: FormElement
}

export type StateSlice = {
  [KEY]: State
}

export const INITIAL_STATE: State = {
  root: {
    type: FieldType.Select
  }
};

export type addFieldProps = {
  path: string,
  field?: FormElement
}
export const addField = createAction<addFieldProps>('formDefinition/addField');

export const formBuilderReducer = createReducer(INITIAL_STATE, (builder => {
  builder.addCase(addField, (state, action) => {
    state = { ...state };
    let path = action.payload.path;

    if (path.startsWith('/')) {
      path = path.substring(1);
    }
    if (path.endsWith('/')) {
      path = path.substring(0, path.length - 1);
    }

    let parts = path.split('/');

    if (parts.length == 0) {
      throw 'addField requires a valid path. Got ' + action.payload.path;
    }
    if ((parts.length - 1) % 2 !== 0) {
      throw 'invalid path. invalid number of parts';
    }
    if (parts[0] !== 'root') {
      throw 'first part of the path must be equal to "root". got ' + parts[0];
    }



    // path can be /root
    //             /root/children/3/


  });
}));

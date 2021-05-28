import {
  CauseType,
  FormDefinition,
  FormDefinitionVersion,
  FormElement, Path, pathFrom,
  Status,
  TranslatedStrings
} from '@core/api-client';
import { createAction, createSlice, PayloadAction } from '@reduxjs/toolkit';

export const KEY = 'formBuilder';

export type State = {
  error: Partial<Status>,
  selectedVersion: string,
  formDefinition: FormDefinition
}

export type StateSlice = {
  formBuilder: State
}

export const INITIAL_STATE: State = {
  error: {
    details: {
      causes: []
    }
  },
  selectedVersion: '',
  formDefinition: {
    kind: 'FormDefinition',
    apiVersion: 'core.nrc.no/v1',
    metadata: {},
    spec: {
      group: '',
      names: {
        kind: '',
        plural: '',
        singular: ''
      },
      versions: [
        {
          name: '',
          storage: true,
          served: true,
          schema: {
            formSchema: {
              root: {}
            }
          }
        }
      ]
    }
  }
};

type setValuePayload = {
  path: string
  value: any
}
export const setValue = createAction<setValuePayload>('formDefinitions/setValue');
const handleSetValue = (state: State, action: PayloadAction<setValuePayload>) => {
  const obj = state.formDefinition;
  const path = pathFrom(action.payload.path)
  path.getValue()
};


type setFormDefinitionPayload = {
  formDefinition: FormDefinition
}
/**
 * Adds a new field to the state at the given path
 * Path is jsonPointer
 * Eg.
 * /root/children/3/children/2
 */
export const setFormDefinition = createAction<setFormDefinitionPayload>('formDefinitions/setState');
const handleSetFormDefinition = (state: State, action: PayloadAction<setFormDefinitionPayload>) => {
  state.formDefinition = action.payload.formDefinition;
};


type addFormElementPayload = {
  path: string,
  field?: FormElement
}
/**
 * Adds a new field to the state at the given path
 * Path is jsonPointer
 * Eg.
 * /root/children/3/children/2
 */
export const addFormElement = createAction<addFormElementPayload>('formDefinitions/addFormElement');
const handleAddFormElement = (state: State, action: PayloadAction<addFormElementPayload>) => {
  const element = findElement(state, action.payload.path);
  if (!element.children) {
    element.children = [];
  }
  element.children.push(action.payload.field);
};


type removeFormElementPayload = {
  path: string
}
/**
 * Removes a field from the state at the given path
 * Path is jsonPointer
 * Eg. /root/children/3
 */
export const removeFormElement = createAction<removeFormElementPayload>('formDefinitions/removeFormElement');
const handleRemoveField = (state: State, action: PayloadAction<removeFormElementPayload>) => {
  const path = clearSlashes(action.payload.path);
  const parent = findParentOf(state, path);
  const parts = path.split('/');
  const idx = parseInt(parts[parts.length - 1]);
  parent.children.splice(idx, 1);
};


type replaceFormElementPayload = {
  path: string,
  field: FormElement
}
/**
 * Replaces a field from the state at the given path
 * Path is jsonPointer
 * eg.: /root, /root/children/3
 */
export const replaceFormElement = createAction<replaceFormElementPayload>('formDefinitions/replaceFormElement');
const handleReplaceFormElement = (state: State, action: PayloadAction<replaceFormElementPayload>) => {
  const path = clearSlashes(action.payload.path);
  const version = findSelectedVersion(state);
  if (path === 'root') {
    version.schema.formSchema.root = action.payload.field;
    return;
  }
  const parts = path.split('/');
  const parent = findParentOf(state, action.payload.path);
  const idx = parseInt(parts[parts.length - 1]);
  parent.children[idx] = action.payload.field;
};


type patchFormElementPayload = {
  path: string,
  field: Partial<FormElement>
}
/**
 * Patches a field from the state at the given path
 * Path is jsonPointer
 * eg.: /root, /root/children/3
 */
export const patchFormElement = createAction<patchFormElementPayload>('formDefinitions/patchFormElement');
const handlePatchFormElement = (state: State, action: PayloadAction<patchFormElementPayload>) => {
  const path = clearSlashes(action.payload.path);
  const root = findCurrentVersionRoot(state);
  let element: FormElement;
  if (path === 'root') {
    element = root;
  } else {
    element = findElement(state, path);
  }
  for (let key of Object.keys(action.payload.field)) {
    if (action.payload.field.hasOwnProperty(key)) {
      element[key] = action.payload.field[key];
    }
  }
};

export type TranslationType = 'label' | 'tooltip'

type setTranslationPayload = {
  path: string,
  locale: string,
  type: TranslationType,
  value: string
}

/**
 * Sets/Adds a translation
 * Path is jsonPointer
 * eg.: /root, /root/children/3
 */
export const setTranslation = createAction<setTranslationPayload>('formDefinitions/setTranslation');
const handleSetTranslation = (state: State, action: PayloadAction<setTranslationPayload>) => {
  const element = findElement(state, action.payload.path);
  let translatedStrings: TranslatedStrings;
  if (action.payload.type === 'label') {
    if (!element.label) {
      element.label = [];
    }
    translatedStrings = element.label;
  } else if (action.payload.type === 'tooltip') {
    if (!element.tooltip) {
      element.tooltip = [];
    }
    translatedStrings = element.tooltip;
  }

  let currentTranslation = translatedStrings.find(t => t.locale === action.payload.locale);
  if (!currentTranslation) {
    currentTranslation = {
      locale: action.payload.locale,
      value: action.payload.value
    };
    translatedStrings.push(currentTranslation);
  } else {
    currentTranslation.value = action.payload.value;
  }
};


type removeTranslationPayload = {
  path: string,
  locale: string,
  type: TranslationType,
}


/**
 * Removes a translation
 * Path is jsonPointer
 * eg.: /root, /root/children/3
 */
export const removeTranslation = createAction<removeTranslationPayload>('formDefinitions/removeTranslation');
const handleRemoveTranslation = (state, action: PayloadAction<removeTranslationPayload>) => {
  const element = findElement(state, action.payload.path);
  let translatedStrings: TranslatedStrings;
  if (action.payload.type === 'label') {
    if (!element.label) {
      return;
    }
    translatedStrings = element.label;
  } else if (action.payload.type === 'tooltip') {
    if (!element.tooltip) {
      return;
    }
    translatedStrings = element.tooltip;
  }

  let currentTranslationIdx = translatedStrings.findIndex(t => t.locale === action.payload.locale);

  if (currentTranslationIdx !== -1) {
    translatedStrings.splice(currentTranslationIdx, 1);
  }

};

/**
 * The reducer for the formBuilder
 */
export const formBuilderSlice = createSlice({
  name: KEY,
  initialState: INITIAL_STATE,
  reducers: {},
  extraReducers: builder => {
    builder.addCase(addFormElement, handleAddFormElement);
    builder.addCase(removeFormElement, handleRemoveField);
    builder.addCase(replaceFormElement, handleReplaceFormElement);
    builder.addCase(patchFormElement, handlePatchFormElement);
    builder.addCase(setTranslation, handleSetTranslation);
    builder.addCase(removeTranslation, handleRemoveTranslation);
    builder.addCase(setFormDefinition, handleSetFormDefinition);
  }
});

export const reducer = formBuilderSlice.reducer;


const validate = (state: State): void => {
  state.error = undefined;
  const status: Partial<Status> = {};
  const formDefinition = state.formDefinition;
  if (!formDefinition.kind) {
    status.details.causes.push({
      type: CauseType.FieldValueInvalid,
      field: 'kind',
      message: 'Kind is required'
    });
  }
};

/**
 * Finds the parent of an element at the given path
 * @param state The current reducer state
 * @param path The path to the child
 */
const findParentOf = (state: State, path: string): FormElement => {
  path = clearSlashes(path);
  let parts = path.split('/');
  parts.splice(parts.length - 2, 2);
  return findElement(state, parts.join('/'));
};

/**
 * Finds an element at the given path
 * @param state The reducer state
 * @param path The path of the element
 */
const findElement = (state: State, path: string): FormElement => {

  path = clearSlashes(path);

  const version = findSelectedVersion(state);
  const root = version.schema.formSchema.root;

  if (path === 'root') {
    return root;
  }

  // Path cannot be empty
  assertNotEmpty(path, 'path cannot be empty');

  let parts = path.split('/');

  // The number of parts must satisfy (n - 1) % 2 == 0
  if ((parts.length - 1) % 2 !== 0) {
    throw 'invalid path. invalid number of parts';
  }

  // The first part must be "root"
  if (parts[0] !== 'root') {
    throw 'first part of the path must be equal to "root", got "' + parts[0] + '"';
  }

  // Walk through the fields to find the target
  let currentField = root;
  parts.splice(0, 1);
  while (parts.length > 0) {
    const currentPart = parts[0];
    console.log(currentPart);
    // @ts-ignore
    if (currentPart === 'children') {
      const idx = parts[1];
      parts.splice(0, 2);
      currentField = currentField.children[parseInt(idx)];
    } else {
      throw 'unexpected part name "' + currentPart + '"';
    }
  }

  return currentField;

};

const findCurrentVersionRoot = (state: State): FormElement => {
  const version = findSelectedVersion(state);
  return version?.schema?.formSchema?.root;
};

const findSelectedVersion = (state: State): FormDefinitionVersion => {
  return findVersion(state, state.selectedVersion);
};

const findVersion = (state: State, name: string): FormDefinitionVersion => {
  return state.formDefinition.spec.versions.find(v => v.name === name);
};

/**
 * Asserts that the given string is not empty, or throws
 * @param str the string
 * @param message the error message to throw if the string is empty
 */
const assertNotEmpty = (str: string, message: string) => {
  if (!str) {
    throw message;
  }
};

/**
 * Removes leading and trailing slashes from a string
 * @param str The string
 */
const clearSlashes = (str: string): string => {
  // Remove the leading slash
  if (str.startsWith('/')) {
    str = str.substring(1);
  }

  // Remove the trailing slash
  if (str.endsWith('/')) {
    str = str.substring(0, str.length - 1);
  }
  return str;
};

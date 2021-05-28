import { FormElement, TranslatedStrings } from '@core/api-client';
import { createAction, createSlice, PayloadAction } from '@reduxjs/toolkit';

export const KEY = 'formBuilder';

export type State = {
  root: FormElement
}

export type StateSlice = {
  formBuilder: State
}

export const INITIAL_STATE: State = {
  root: {
    type: 'float'
  }
};


/**
 * Adds a new field to the state at the given path
 * Path is jsonPointer
 * Eg.
 * /root/children/3/children/2
 */
export const setState = createAction<{
  state: State
}>('formDefinitions/setState');
const handleSetState = (state, action) => {
  return action.payload.state;
};

/**
 * Adds a new field to the state at the given path
 * Path is jsonPointer
 * Eg.
 * /root/children/3/children/2
 */
export const addField = createAction<{
  path: string,
  field?: FormElement
}>('formDefinitions/addField');

const handleAddField = (state, action) => {
  const element = findElement(state, action.payload.path);
  if (!element.children) {
    element.children = [];
  }
  element.children.push(action.payload.field);
};


/**
 * Removes a field from the state at the given path
 * Path is jsonPointer
 * Eg. /root/children/3
 */
export const removeField = createAction<{
  path: string
}>('formDefinitions/removeField');
const handleRemoveField = (state, action) => {
  const path = clearSlashes(action.payload.path);
  const parent = findParentOf(state, path);
  const parts = path.split('/');
  const idx = parseInt(parts[parts.length - 1]);
  parent.children.splice(idx, 1);
};


/**
 * Replaces a field from the state at the given path
 * Path is jsonPointer
 * eg.: /root, /root/children/3
 */
export const replaceField = createAction<{
  path: string,
  field: FormElement
}>('formDefinitions/replaceField');
const handleReplaceField = (state, action) => {
  const path = clearSlashes(action.payload.path);
  if (path === 'root') {
    state.root = action.payload.field;
    return;
  }
  const parts = path.split('/');
  const parent = findParentOf(state, action.payload.path);
  const idx = parseInt(parts[parts.length - 1]);
  parent.children[idx] = action.payload.field;
};


/**
 * Patches a field from the state at the given path
 * Path is jsonPointer
 * eg.: /root, /root/children/3
 */
export const patchField = createAction<{
  path: string,
  field: Partial<FormElement>
}>('formDefinitions/patchField');
const handlePatchField = (state, action) => {
  const path = clearSlashes(action.payload.path);
  let element: FormElement;
  if (path === 'root') {
    element = state.root;
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
const handleSetTranslation = (state, action: PayloadAction<setTranslationPayload>) => {
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
    builder.addCase(addField, handleAddField);
    builder.addCase(removeField, handleRemoveField);
    builder.addCase(replaceField, handleReplaceField);
    builder.addCase(patchField, handlePatchField);
    builder.addCase(setTranslation, handleSetTranslation);
    builder.addCase(removeTranslation, handleRemoveTranslation);
    builder.addCase(setState, handleSetState);
  }
});

export const reducer = formBuilderSlice.reducer;


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

  if (path === 'root') {
    return state.root;
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
  let currentField = state.root;
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

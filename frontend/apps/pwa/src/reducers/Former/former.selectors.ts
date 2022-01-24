import {
  Database,
  FieldDefinition,
  FieldType,
  Folder,
  FormDefinition,
} from 'core-api-client';

import { RootState } from '../../app/store';
import { databaseGlobalSelectors } from '../database';
import { folderGlobalSelectors } from '../folder';

import { Form, FormerState, FormField } from './types';
import { adapter } from './former.adapter';

const selectors = adapter.getSelectors();
const globalSelectors = adapter.getSelectors<RootState>(
  (state) => state.former,
);

const selectFieldForm = (
  state: FormerState,
  fieldId: string,
): { form: Form; field: FormField } | undefined => {
  const allForms: Form[] = [];
  for (const id of state.ids) {
    const entity = state.entities[id];
    if (entity) {
      allForms.push(entity);
    }
  }
  let field: FormField | undefined;
  let form: Form | undefined;
  for (const candidateForm of allForms) {
    for (const formField of candidateForm.fields) {
      if (formField.id === fieldId) {
        form = candidateForm;
        field = formField;
      }
    }
  }
  if (!field || !form) {
    return;
  }
  return { form, field };
};

const selectFormFields = (state: FormerState, formId: string): FormField[] => {
  const form = selectors.selectById(state, formId);
  if (!form) {
    return [];
  }
  return form.fields;
};

const selectCurrentForm = (state: FormerState): Form | undefined => {
  return selectors.selectById(state, state.selectedFormId);
};

const selectCurrentFormFields = (state: FormerState): FormField[] => {
  const form = selectors.selectById(state, state.selectedFormId);
  if (!form) {
    return [];
  }
  return form.fields;
};

const selectCurrentField = (state: FormerState): FormField | undefined => {
  if (!state.selectedFieldId) {
    return undefined;
  }
  const formField = selectFieldForm(state, state.selectedFieldId);
  if (!formField) {
    return undefined;
  }
  return formField.field;
};

const selectOwnerForm = (
  state: FormerState,
  subFormId: string,
): Form | undefined => {
  const allForms = selectors.selectAll(state);
  for (const form of allForms) {
    for (const field of form.fields) {
      if (field.subFormId === subFormId) {
        return form;
      }
    }
  }
  return undefined;
};

const selectCurrentFormOwner = (state: FormerState): Form | undefined => {
  return selectOwnerForm(state, state.selectedFormId);
};

const selectIsSubForm = (state: FormerState, formId: string): boolean => {
  return !!selectOwnerForm(state, formId);
};

const selectIsRootForm = (state: FormerState, formId: string): boolean => {
  return !selectIsSubForm(state, formId);
};

const selectDatabase = (state: RootState): Database | undefined => {
  if (!state.former.selectedDatabaseId) {
    return undefined;
  }
  return databaseGlobalSelectors.selectById(
    state,
    state.former.selectedDatabaseId,
  );
};

const selectFolder = (state: RootState): Folder | undefined => {
  if (!state.former.selectedFolderId) {
    return undefined;
  }
  return folderGlobalSelectors.selectById(state, state.former.selectedFolderId);
};

const mapFields = (
  state: FormerState,
  fields: FormField[],
): FieldDefinition[] => {
  const result: FieldDefinition[] = [];

  if (!fields) {
    return [];
  }

  for (const field of fields) {
    let fieldType: FieldType;

    if (field.fieldType === 'text') {
      fieldType = { text: {} };
    } else if (field.fieldType === 'multilineText') {
      fieldType = { multilineText: {} };
    } else if (field.fieldType === 'date') {
      fieldType = { date: {} };
    } else if (field.fieldType === 'month') {
      fieldType = { month: {} };
    } else if (field.fieldType === 'week') {
      fieldType = { week: {} };
    } else if (field.fieldType === 'quantity') {
      fieldType = { quantity: {} };
    } else if (field.fieldType === 'checkbox') {
      fieldType = { checkbox: {} };
    } else if (field.fieldType === 'singleSelect') {
      fieldType = { singleSelect: { options: field.options } };
    } else if (field.fieldType === 'multiSelect') {
      fieldType = { multiSelect: { options: field.options } };
    } else if (field.fieldType === 'reference') {
      if (!field.referencedDatabaseId) {
        throw new Error(
          `field with id ${field.id} does not have referenced database id`,
        );
      }
      if (!field.referencedFormId) {
        throw new Error(
          `field with id ${field.id} does not have referenced form id`,
        );
      }
      fieldType = {
        reference: {
          databaseId: field.referencedDatabaseId,
          formId: field.referencedFormId,
        },
      };
    } else if (field.fieldType === 'subform') {
      if (!field.subFormId) {
        throw new Error(
          `subform field with id ${field.id} does not have subFormId`,
        );
      }
      const subForm = selectors.selectById(state, field.subFormId);
      if (!subForm) {
        throw new Error(`subform with id ${field.subFormId} not found`);
      }
      fieldType = {
        subForm: {
          fields: mapFields(state, subForm.fields),
        },
      };
    } else {
      throw new Error(
        `invalid field type form field ${field.id}: ${field.fieldType}`,
      );
    }

    result.push({
      fieldType,
      id: '',
      description: field.description,
      name: field.name,
      required: field.required,
      code: field.code,
      key: field.key,
    });
  }

  return result;
};

function selectFormDefinition(
  databaseId: string | undefined,
  folderId: string | undefined,
): (state: FormerState) => Partial<FormDefinition> | undefined {
  return (state) => {
    try {
      const allForms = selectors.selectAll(state);
      const rootForm = allForms.find((e) => e.isRootForm);
      if (!rootForm) {
        return;
      }
      return {
        databaseId,
        folderId,
        name: rootForm.name,
        id: '',
        fields: mapFields(state, rootForm.fields),
        code: '',
      };
    } catch (err) {
      return undefined;
    }
  };
}

export const formerSelectors = {
  ...selectors,
  selectFieldForm,
  selectFormFields,
  selectCurrentForm,
  selectCurrentFormFields,
  selectCurrentField,
  selectIsSubForm,
  selectIsRootForm,
  selectOwnerForm,
  selectCurrentFormOwner,
  selectFormDefinition,
};

export const formerGlobalSelectors = {
  globalSelectors,
  selectDatabase,
  selectFolder,
  selectFieldForm: (state: RootState, fieldId: string) =>
    selectFieldForm(state.former, fieldId),
  selectFormFields: (state: RootState, formId: string) =>
    selectFormFields(state.former, formId),
  selectCurrentForm: (state: RootState) => selectCurrentForm(state.former),
  selectCurrentFormFields: (state: RootState) =>
    selectCurrentFormFields(state.former),
  selectCurrentField: (state: RootState) => selectCurrentField(state.former),
  selectIsSubForm: (state: RootState, formId: string) =>
    selectIsSubForm(state.former, formId),
  selectIsRootForm: (state: RootState, formId: string) =>
    selectIsRootForm(state.former, formId),
  selectOwnerForm: (state: RootState, formId: string) =>
    selectOwnerForm(state.former, formId),
  selectCurrentFormOwner: (state: RootState) =>
    selectCurrentFormOwner(state.former),
  selectFormDefinition:
    (databaseId: string | undefined, folderId: string | undefined) =>
    (state: RootState) => {
      return selectFormDefinition(databaseId, folderId)(state.former);
    },
};

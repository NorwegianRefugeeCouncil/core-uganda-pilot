import {
  createAsyncThunk,
  createEntityAdapter,
  createSlice,
  EntityState,
  PayloadAction,
} from '@reduxjs/toolkit';
import { v4 as uuidv4 } from 'uuid';
import {
  Database,
  FieldDefinition,
  FieldKind,
  FieldType,
  Folder,
  FormDefinition,
  SelectOption,
} from 'core-api-client';

import { RootState } from '../app/store';
import client from '../app/client';

import { databaseGlobalSelectors } from './database';
import { folderGlobalSelectors } from './folder';

export interface FormField {
  id: string;
  type: FieldKind;
  options: SelectOption[];
  required: boolean;
  key: boolean;
  name: string;
  description: string;
  code: string;
  subFormId: string | undefined;
  referencedDatabaseId: string | undefined;
  referencedFormId: string | undefined;
}

export interface Form {
  // name of the form
  name: string;
  // the unique id of the form
  formId: string;
  // records the record values
  fields: FormField[];

  isRootForm: boolean;
}

const adapter = createEntityAdapter<Form>({
  // Assume IDs are stored in a field other than `book.id`
  selectId: (folder) => folder.formId,
  // Keep the "all IDs" array sorted based on book titles
  sortComparer: (a, b) => a.formId.localeCompare(b.formId),
});

export const postForm = createAsyncThunk<
  FormDefinition,
  Partial<FormDefinition>
>('former/createForm', async (arg, thunkAPI) => {
  const resp = await client.createForm({ object: arg });
  if (resp.success) {
    return resp.response as FormDefinition;
  }
  console.log('RESPONSE', resp.error);
  throw resp.error;
});

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
  if (!fields) {
    return [];
  }

  const result = fields.map<FieldDefinition>((field) => {
    let fieldType: FieldType;

    if (field.type === 'text') {
      fieldType = { text: {} };
    } else if (field.type === 'multilineText') {
      fieldType = { multilineText: {} };
    } else if (field.type === 'date') {
      fieldType = { date: {} };
    } else if (field.type === 'month') {
      fieldType = { month: {} };
    } else if (field.type === 'week') {
      fieldType = { week: {} };
    } else if (field.type === 'quantity') {
      fieldType = { quantity: {} };
    } else if (field.type === 'checkbox') {
      fieldType = { checkbox: {} };
    } else if (field.type === 'singleSelect') {
      fieldType = { singleSelect: { options: field.options } };
    } else if (field.type === 'multiSelect') {
      fieldType = { multiSelect: { options: field.options } };
    } else if (field.type === 'reference') {
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
    } else if (field.type === 'subform') {
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
        `invalid field type form field ${field.id}: ${field.type}`,
      );
    }

    return {
      fieldType,
      id: '',
      description: field.description,
      name: field.name,
      required: field.required,
      code: field.code,
      key: field.key,
    };
  });

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
        return undefined;
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

interface FormerState extends EntityState<Form> {
  selectedFormId: string;
  selectedFieldId: string | undefined;
  selectedDatabaseId: string | undefined;
  selectedFolderId: string | undefined;
  savePending: boolean;
  saveSuccess: boolean;
  saveError: any;
}

export const former = createSlice({
  name: 'recorder',
  initialState: {
    ...adapter.getInitialState(),
    selectedFormId: '',
    selectedFieldId: undefined as string | undefined,
    selectedDatabaseId: undefined as string | undefined,
    selectedFolderId: undefined as string | undefined,
    savePending: false,
    saveSuccess: false,
    saveId: undefined as string | undefined,
    saveError: undefined,
  },
  reducers: {
    reset(state) {
      state.entities = {};
      state.ids = [];
      const formId = uuidv4();
      state.selectedFormId = formId;
      state.selectedFieldId = undefined;
      adapter.addOne(state, {
        formId,
        fields: [],
        name: '',
        isRootForm: true,
      });
    },
    setDatabase(state, action: PayloadAction<{ databaseId: string }>) {
      state.selectedDatabaseId = action.payload.databaseId;
    },
    setFolder(state, action: PayloadAction<{ folderId: string }>) {
      state.selectedFolderId = action.payload.folderId;
    },
    setFormName(
      state,
      action: PayloadAction<{ formId: string; formName: string }>,
    ) {
      const { formId, formName } = action.payload;
      const form = state.entities[formId];
      if (form) {
        form.name = formName;
      }
    },
    setFieldRequired(
      state,
      action: PayloadAction<{ fieldId: string; required: boolean }>,
    ) {
      const fieldForm = selectFieldForm(state, action.payload.fieldId);
      if (!fieldForm) {
        return;
      }
      fieldForm.field.required = action.payload.required;
    },
    setFieldIsKey(
      state,
      action: PayloadAction<{ fieldId: string; isKey: boolean }>,
    ) {
      const fieldForm = selectFieldForm(state, action.payload.fieldId);
      if (!fieldForm) {
        return;
      }
      fieldForm.field.required = true;
      fieldForm.field.key = action.payload.isKey;
    },
    setFieldName(
      state,
      action: PayloadAction<{ fieldId: string; name: string }>,
    ) {
      const fieldForm = selectFieldForm(state, action.payload.fieldId);
      if (!fieldForm) {
        return;
      }
      fieldForm.field.name = action.payload.name;
    },
    setFieldDescription(
      state,
      action: PayloadAction<{ fieldId: string; description: string }>,
    ) {
      const fieldForm = selectFieldForm(state, action.payload.fieldId);
      if (!fieldForm) {
        return;
      }
      fieldForm.field.description = action.payload.description;
    },
    setFieldOption(
      state,
      action: PayloadAction<{ fieldId: string; i: number; value: string }>,
    ) {
      const fieldForm = selectFieldForm(state, action.payload.fieldId);
      const { i, value } = action.payload;
      if (!fieldForm || !fieldForm.field.options) return;

      fieldForm.field.options[i] = {
        id: fieldForm.field.options[i].id,
        name: value,
      };
    },
    addOption(state, action: PayloadAction<{ fieldId: string }>) {
      const fieldForm = selectFieldForm(state, action.payload.fieldId);
      if (!fieldForm) return;

      fieldForm.field.options = [
        ...fieldForm.field.options,
        {
          id: uuidv4(),
          name: '',
        },
      ];
    },
    removeOption(state, action: PayloadAction<{ fieldId: string; i: number }>) {
      const fieldForm = selectFieldForm(state, action.payload.fieldId);
      if (!fieldForm || !fieldForm.field.options) return;

      const { i } = action.payload;
      fieldForm.field.options = fieldForm.field.options
        .slice(0, i)
        .concat(fieldForm.field.options.slice(i + 1));
    },
    setFieldCode(
      state,
      action: PayloadAction<{ fieldId: string; code: string }>,
    ) {
      const fieldForm = selectFieldForm(state, action.payload.fieldId);
      if (!fieldForm) {
        return;
      }
      fieldForm.field.code = action.payload.code;
    },
    setFieldReferencedDatabaseId(
      state,
      action: PayloadAction<{ fieldId: string; databaseId: string }>,
    ) {
      const fieldForm = selectFieldForm(state, action.payload.fieldId);
      if (!fieldForm) {
        return;
      }
      fieldForm.field.referencedDatabaseId = action.payload.databaseId;
    },
    setFieldReferencedFormId(
      state,
      action: PayloadAction<{ fieldId: string; formId: string }>,
    ) {
      const fieldForm = selectFieldForm(state, action.payload.fieldId);
      if (!fieldForm) {
        return;
      }
      fieldForm.field.referencedFormId = action.payload.formId;
    },
    cancelFieldChanges(state, action: PayloadAction<{ fieldId: string }>) {
      const { fieldId } = action.payload;
      let form: Form | undefined;
      let field: FormField | undefined;
      for (const formId in state.entities) {
        if (!state.entities.hasOwnProperty(formId)) {
          continue;
        }
        const candidateForm = state.entities[formId];
        if (!candidateForm) {
          continue;
        }
        for (const formField of candidateForm.fields) {
          if (formField.id === fieldId) {
            form = candidateForm;
            field = formField;
            break;
          }
        }
      }

      if (!form) {
        return;
      }
      if (!field) {
        return;
      }

      form.fields = form.fields.filter((f) => f.id !== field?.id);

      if (state.selectedFieldId === action.payload.fieldId) {
        state.selectedFieldId = undefined;
      }
    },
    addField(
      state,
      action: PayloadAction<{
        formId: string;
        kind: FieldKind;
        referencedDatabaseId?: string;
        referencedFormId?: string;
      }>,
    ) {
      const { formId, kind, referencedDatabaseId, referencedFormId } =
        action.payload;
      const form = formerSelectors.selectById(state, formId);
      if (!form) {
        return;
      }
      const fieldId = uuidv4();
      let subFormId: string | undefined;
      if (kind === FieldKind.SubForm) {
        subFormId = uuidv4();
        const subForm: Form = {
          formId: subFormId,
          fields: [],
          name: '',
          isRootForm: false,
        };
        adapter.addOne(state, subForm);
      }
      const newField: FormField = {
        id: fieldId,
        key: false,
        name: '',
        required: false,
        type: kind,
        options: [],
        subFormId,
        code: '',
        description: '',
        referencedDatabaseId,
        referencedFormId,
      };
      state.entities[form.formId] = {
        ...form,
        fields: [...form.fields, newField],
      };
      state.selectedFieldId = fieldId;
    },
    addSubForm(state, action: PayloadAction<{ ownerFieldId: string }>) {
      const newForm: Form = {
        formId: uuidv4(),
        fields: [],
        name: '',
        isRootForm: false,
      };
      adapter.addOne(state, newForm);
    },
    openSubForm(state, action: PayloadAction<{ fieldId: string }>) {
      const fieldForm = selectFieldForm(state, action.payload.fieldId);
      if (!fieldForm) {
        return;
      }
      if (!fieldForm.field.subFormId) {
        return;
      }
      const subForm = formerSelectors.selectById(
        state,
        fieldForm.field.subFormId,
      );
      if (!subForm) {
        return;
      }
      state.selectedFormId = subForm.formId;
    },
    selectForm(state, action: PayloadAction<{ formId: string }>) {
      const { formId } = action.payload;
      const form = formerSelectors.selectById(state, formId);
      if (!form) {
        return;
      }
      state.selectedFormId = form.formId;
    },
    saveForm(state) {
      const ownerForm = selectOwnerForm(state, state.selectedFormId);
      if (ownerForm) {
        state.selectedFieldId = undefined;
        state.selectedFormId = ownerForm.formId;
      } else {
      }
    },
    selectField(state, action: PayloadAction<{ fieldId: string | undefined }>) {
      if (!action.payload.fieldId) {
        state.selectedFieldId = undefined;
        return;
      }
      const fieldForm = selectFieldForm(state, action.payload.fieldId);
      if (!fieldForm) {
        return;
      }
      state.selectedFormId = fieldForm.form.formId;
      state.selectedFieldId = fieldForm.field.id;
    },
    addOne: adapter.addOne,
    addMany: adapter.addMany,
    removeAll: adapter.removeAll,
    removeMany: adapter.removeMany,
    removeOne: adapter.removeOne,
    updateMany: adapter.updateMany,
    updateOne: adapter.updateOne,
    upsertOne: adapter.upsertOne,
    upsertMany: adapter.upsertMany,
    setOne: adapter.setOne,
    setMany: adapter.setMany,
    setAll: adapter.setAll,
  },
  extraReducers: (builder) => {
    builder.addCase(postForm.pending, (state) => {
      state.savePending = true;
    });
    builder.addCase(postForm.rejected, (state, payload) => {
      state.savePending = false;
      state.saveError = payload.payload as any;
      state.saveSuccess = false;
    });
    builder.addCase(postForm.fulfilled, (state, payload) => {
      state.savePending = false;
      state.saveError = undefined;
      state.saveSuccess = true;
      state.saveId = payload.payload.id;
    });
  },
});

export const formerActions = former.actions;
export default former.reducer;

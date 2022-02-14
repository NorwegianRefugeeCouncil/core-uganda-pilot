import { v4 as uuidv4 } from 'uuid';
import {
  ActionReducerMapBuilder,
  createAsyncThunk,
  PayloadAction,
} from '@reduxjs/toolkit';
import { FieldKind, FormDefinition, FormType } from 'core-api-client';

import client from '../../app/client';

import { adapter } from './former.adapter';
import { formerSelectors } from './former.selectors';
import { Form, FormerState, FormField } from './types';

export const reducers = {
  reset(state: FormerState) {
    state.entities = {};
    state.ids = [];
    const formId = uuidv4();
    state.selectedFormId = formId;
    state.selectedFieldId = undefined;
    adapter.addOne(state, {
      formId,
      fields: [],
      name: '',
      formType: FormType.DefaultFormType,
      isRootForm: true,
    });
  },
  setDatabase(
    state: FormerState,
    action: PayloadAction<{ databaseId: string }>,
  ) {
    state.selectedDatabaseId = action.payload.databaseId;
  },
  setFolder(state: FormerState, action: PayloadAction<{ folderId: string }>) {
    state.selectedFolderId = action.payload.folderId;
  },
  setFormName(
    state: FormerState,
    action: PayloadAction<{ formId: string; formName: string }>,
  ) {
    const { formId, formName } = action.payload;
    const form = state.entities[formId];
    if (form) {
      form.name = formName;
    }
  },
  setFieldRequired(
    state: FormerState,
    action: PayloadAction<{ fieldId: string; required: boolean }>,
  ) {
    const fieldForm = formerSelectors.selectFieldForm(
      state,
      action.payload.fieldId,
    );
    if (!fieldForm) {
      return;
    }
    fieldForm.field.required = action.payload.required;
  },
  setFieldIsKey(
    state: FormerState,
    action: PayloadAction<{ fieldId: string; isKey: boolean }>,
  ) {
    const fieldForm = formerSelectors.selectFieldForm(
      state,
      action.payload.fieldId,
    );
    if (!fieldForm) {
      return;
    }
    fieldForm.field.required = true;
    fieldForm.field.key = action.payload.isKey;
  },
  setFieldName(
    state: FormerState,
    action: PayloadAction<{ fieldId: string; name: string }>,
  ) {
    const fieldForm = formerSelectors.selectFieldForm(
      state,
      action.payload.fieldId,
    );
    if (!fieldForm) {
      return;
    }
    fieldForm.field.name = action.payload.name;
  },
  setFieldDescription(
    state: FormerState,
    action: PayloadAction<{ fieldId: string; description: string }>,
  ) {
    const fieldForm = formerSelectors.selectFieldForm(
      state,
      action.payload.fieldId,
    );
    if (!fieldForm) {
      return;
    }
    fieldForm.field.description = action.payload.description;
  },
  setFieldOption(
    state: FormerState,
    action: PayloadAction<{ fieldId: string; i: number; value: string }>,
  ) {
    const fieldForm = formerSelectors.selectFieldForm(
      state,
      action.payload.fieldId,
    );
    const { i, value } = action.payload;
    if (!fieldForm || !fieldForm.field.options) return;

    fieldForm.field.options[i] = {
      id: fieldForm.field.options[i].id,
      name: value,
    };
  },
  addOption(state: FormerState, action: PayloadAction<{ fieldId: string }>) {
    const fieldForm = formerSelectors.selectFieldForm(
      state,
      action.payload.fieldId,
    );
    if (!fieldForm) return;

    fieldForm.field.options = [
      ...fieldForm.field.options,
      {
        id: uuidv4(),
        name: '',
      },
    ];
  },
  removeOption(
    state: FormerState,
    action: PayloadAction<{ fieldId: string; i: number }>,
  ) {
    const fieldForm = formerSelectors.selectFieldForm(
      state,
      action.payload.fieldId,
    );
    if (!fieldForm || !fieldForm.field.options) return;

    const { i } = action.payload;
    fieldForm.field.options = fieldForm.field.options
      .slice(0, i)
      .concat(fieldForm.field.options.slice(i + 1));
  },
  setFieldCode(
    state: FormerState,
    action: PayloadAction<{ fieldId: string; code: string }>,
  ) {
    const fieldForm = formerSelectors.selectFieldForm(
      state,
      action.payload.fieldId,
    );
    if (!fieldForm) {
      return;
    }
    fieldForm.field.code = action.payload.code;
  },
  setFieldReferencedDatabaseId(
    state: FormerState,
    action: PayloadAction<{ fieldId: string; databaseId: string }>,
  ) {
    const fieldForm = formerSelectors.selectFieldForm(
      state,
      action.payload.fieldId,
    );
    if (!fieldForm) {
      return;
    }
    fieldForm.field.referencedDatabaseId = action.payload.databaseId;
  },
  setFieldReferencedFormId(
    state: FormerState,
    action: PayloadAction<{ fieldId: string; formId: string }>,
  ) {
    const fieldForm = formerSelectors.selectFieldForm(
      state,
      action.payload.fieldId,
    );
    if (!fieldForm) {
      return;
    }
    fieldForm.field.referencedFormId = action.payload.formId;
  },
  cancelFieldChanges(
    state: FormerState,
    action: PayloadAction<{ fieldId: string }>,
  ) {
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
    state: FormerState,
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
        formType: FormType.DefaultFormType,
        isRootForm: false,
      };
      adapter.addOne(state, subForm);
    }
    const newField: FormField = {
      id: fieldId,
      key: false,
      name: '',
      required: false,
      fieldType: kind,
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
  addSubForm(
    state: FormerState,
    action: PayloadAction<{ ownerFieldId: string }>,
  ) {
    const newForm: Form = {
      formId: uuidv4(),
      fields: [],
      name: '',
      formType: FormType.DefaultFormType,
      isRootForm: false,
    };
    adapter.addOne(state, newForm);
  },
  openSubForm(state: FormerState, action: PayloadAction<{ fieldId: string }>) {
    const fieldForm = formerSelectors.selectFieldForm(
      state,
      action.payload.fieldId,
    );
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
    state.selectedFieldId = undefined;
  },
  selectForm(state: FormerState, action: PayloadAction<{ formId: string }>) {
    const { formId } = action.payload;
    const form = formerSelectors.selectById(state, formId);
    if (!form) {
      return;
    }
    state.selectedFormId = form.formId;
  },
  saveForm(state: FormerState) {
    const ownerForm = formerSelectors.selectOwnerForm(
      state,
      state.selectedFormId,
    );
    if (ownerForm) {
      state.selectedFieldId = undefined;
      state.selectedFormId = ownerForm.formId;
    } else {
    }
  },
  selectField(
    state: FormerState,
    action: PayloadAction<{ fieldId: string | undefined }>,
  ) {
    if (!action.payload.fieldId) {
      state.selectedFieldId = undefined;
      return;
    }
    const fieldForm = formerSelectors.selectFieldForm(
      state,
      action.payload.fieldId,
    );
    if (!fieldForm) {
      return;
    }
    state.selectedFormId = fieldForm.form.formId;
    state.selectedFieldId = fieldForm.field.id;
  },
  setFormType(
    state: FormerState,
    action: PayloadAction<{ formId: string; formType: FormType }>,
  ) {
    const { formId, formType } = action.payload;
    const form = state.entities[formId];
    if (!form || form.formType === formType) return;

    form.formType = formType;
    if (formType === FormType.DefaultFormType) {
      // If switching to default form type, remove recipient reference field
      form.fields = form.fields.filter(
        (f) => f.fieldType === FieldKind.Reference && f.key,
      );
    } else if (formType === FormType.RecipientFormType) {
      // If switching to recipient form type, add recipient reference field
      const newFieldId = uuidv4();
      form.fields = [
        {
          id: newFieldId,
          name: '',
          code: '',
          description: '',
          key: true,
          required: true,
          fieldType: FieldKind.Reference,
          subFormId: undefined,
          referencedDatabaseId: undefined,
          referencedFormId: undefined,
          options: [],
        },
        ...form.fields,
      ];
      state.selectedFieldId = newFieldId;
    }

    state.entities[formId] = form;
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
};

export const postForm = createAsyncThunk<
  FormDefinition,
  Partial<FormDefinition>
>('former/createForm', async (arg, thunkAPI) => {
  const resp = await client.Form.create({ object: arg });
  if (resp.success) {
    return resp.response as FormDefinition;
  }
  return thunkAPI.rejectWithValue(resp?.error);
});

export const extraReducers = (builder: ActionReducerMapBuilder<any>) => {
  builder.addCase(postForm.pending, (state) => {
    state.savePending = true;
    state.saveError = undefined;
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
};

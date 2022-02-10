import { FC, useCallback, useEffect } from 'react';
import { useHistory, useLocation } from 'react-router-dom';
import { FieldKind } from 'core-api-client';
import { useForm, FormProvider } from 'react-hook-form';

import { useAppDispatch, useAppSelector } from '../../app/hooks';
import { fetchDatabases } from '../../reducers/database';
import { fetchFolders } from '../../reducers/folder';
import { fetchForms } from '../../reducers/form';
import former from '../../reducers/Former';
import { formerGlobalSelectors } from '../../reducers/Former/former.selectors';
import { postForm } from '../../reducers/Former/former.reducers';
import { FormField, ValidationForm } from '../../reducers/Former/types';

import { Former } from './Former';
import { customValidation } from './validation';

export const FormerContainer: FC = () => {
  const dispatch = useAppDispatch();
  const history = useHistory();
  const { actions } = former;

  // load data
  useEffect(() => {
    dispatch(actions.reset());
    dispatch(fetchDatabases());
    dispatch(fetchFolders());
    dispatch(fetchForms());
  }, [dispatch]);

  const location = useLocation();

  const form = useAppSelector(formerGlobalSelectors.selectCurrentForm);
  const ownerForm = useAppSelector(
    formerGlobalSelectors.selectCurrentFormOwner,
  );
  const folder = useAppSelector(formerGlobalSelectors.selectFolder);
  const database = useAppSelector(formerGlobalSelectors.selectDatabase);
  const selectedField = useAppSelector(
    formerGlobalSelectors.selectCurrentField,
  );

  const formDefinition = useAppSelector(
    formerGlobalSelectors.selectFormDefinition(database?.id, folder?.id),
  );
  const useFormObject = useForm<ValidationForm>({
    defaultValues: form,
  });
  const {
    clearErrors,
    formState,
    handleSubmit,
    reset,
    resetField,
    setError,
    trigger,
  } = useFormObject;

  useEffect(() => {
    const search = new URLSearchParams(location.search);
    const databaseId = search.get('databaseId');
    if (databaseId) {
      dispatch(actions.setDatabase({ databaseId }));
    }
    const folderId = search.get('folderId');
    if (folderId) {
      dispatch(actions.setFolder({ folderId }));
    }
  }, [dispatch, location]);

  const setFormName = useCallback(
    (formName: string) => {
      if (form) {
        dispatch(actions.setFormName({ formId: form.formId, formName }));
      }
    },
    [dispatch, form],
  );

  const setSelectedField = useCallback(
    (fieldId: string | undefined) => {
      reset();
      clearErrors();
      dispatch(actions.selectField({ fieldId }));
    },
    [dispatch],
  );

  const addField = useCallback(
    (kind: FieldKind) => {
      resetField('selectedField');
      if (form) {
        dispatch(actions.addField({ formId: form.formId, kind }));
        clearErrors('fields');
      }
    },
    [dispatch, form],
  );

  const setFieldOption = useCallback(
    (fieldId: string, i: number, value: string) => {
      dispatch(actions.setFieldOption({ fieldId, i, value }));
    },
    [dispatch],
  );

  const addOption = useCallback(
    (fieldId: string) => {
      clearErrors('selectedField.fieldType.singleSelect.options');
      clearErrors('selectedField.fieldType.multiSelect.options');
      dispatch(actions.addOption({ fieldId }));
    },
    [dispatch],
  );

  const removeOption = useCallback(
    (fieldId: string, i: number) => {
      dispatch(actions.removeOption({ fieldId, i }));
    },
    [dispatch],
  );

  const cancelField = useCallback(
    (fieldId: string) => {
      dispatch(actions.cancelFieldChanges({ fieldId }));
      resetField('selectedField');
    },
    [dispatch],
  );

  const setFieldRequired = useCallback(
    (fieldId: string, required: boolean) => {
      dispatch(actions.setFieldRequired({ fieldId, required }));
    },
    [dispatch],
  );

  const setFieldIsKey = useCallback(
    (fieldId: string, isKey: boolean) => {
      dispatch(actions.setFieldIsKey({ fieldId, isKey }));
    },
    [dispatch],
  );

  const setFieldName = useCallback(
    (fieldId: string, name: string) => {
      dispatch(actions.setFieldName({ fieldId, name }));
    },
    [dispatch],
  );

  const setFieldDescription = useCallback(
    (fieldId: string, description: string) => {
      dispatch(actions.setFieldDescription({ fieldId, description }));
    },
    [dispatch],
  );

  const setFieldReferencedDatabaseId = useCallback(
    (fieldId: string, databaseId: string) => {
      dispatch(actions.setFieldReferencedDatabaseId({ fieldId, databaseId }));
    },
    [dispatch],
  );

  const setFieldReferencedFormId = useCallback(
    (fieldId: string, formId: string) => {
      dispatch(actions.setFieldReferencedFormId({ fieldId, formId }));
    },
    [dispatch],
  );

  const openSubForm = useCallback(
    (fieldId: string) => {
      dispatch(actions.openSubForm({ fieldId }));
    },
    [dispatch],
  );

  const saveField = useCallback(
    async (field: FormField) => {
      const valid = await trigger('selectedField');
      const errors = customValidation.selectedField(field);

      if (errors.length) {
        errors.forEach((error) => {
          setError(error.field, {
            message: error.message,
          });
        });
        return;
      }
      if (valid) {
        dispatch(actions.selectField({ fieldId: undefined }));
      }
    },
    [dispatch],
  );

  const saveForm = useCallback(async () => {
    if (form) {
      const errors = customValidation.form(form);
      if (errors.length) {
        errors.forEach((error) => {
          setError(error.field, {
            message: error.message,
          });
        });
        return;
      }
    }

    if (ownerForm) {
      dispatch(actions.saveForm());
    } else if (formDefinition) {
      try {
        const data = await dispatch(postForm(formDefinition)).unwrap();
        history.push(`/browse/forms/${data.id}`);
      } catch (apiErrors: any) {
        apiErrors?.details?.causes.forEach((error: any) => {
          setError(error.field, { type: error.reason, message: error.message });
        });
      }
    }
  }, [dispatch, formDefinition, ownerForm]);

  if (!form) {
    return <></>;
  }

  return (
    <FormProvider {...useFormObject}>
      <Former
        formId={form.formId}
        formType={form.formType}
        addField={addField}
        addOption={addOption}
        cancelField={(fieldId: string) => cancelField(fieldId)}
        errors={formState.errors}
        fields={form.fields}
        formName={form.name}
        invalid={!formState.isValid && formState.isDirty}
        openSubForm={openSubForm}
        ownerFormName={ownerForm?.name}
        removeOption={removeOption}
        saveField={saveField}
        saveForm={handleSubmit(saveForm)}
        selectedFieldId={selectedField?.id}
        setFieldDescription={setFieldDescription}
        setFieldIsKey={setFieldIsKey}
        setFieldName={setFieldName}
        setFieldOption={setFieldOption}
        setFieldReferencedDatabaseId={setFieldReferencedDatabaseId}
        setFieldReferencedFormId={setFieldReferencedFormId}
        setFieldRequired={setFieldRequired}
        setFormName={setFormName}
        setSelectedField={setSelectedField}
      />
    </FormProvider>
  );
};

import React, { FC, useCallback, useEffect } from 'react';
import { useHistory, useLocation } from 'react-router-dom';
import { FieldKind } from 'core-api-client';
import { useForm } from 'react-hook-form';

import { useAppDispatch, useAppSelector } from '../../app/hooks';
import { fetchDatabases } from '../../reducers/database';
import { fetchFolders } from '../../reducers/folder';
import { fetchForms } from '../../reducers/form';
import former from '../../reducers/Former';
import { formerGlobalSelectors } from '../../reducers/Former/former.selectors';
import { postForm } from '../../reducers/Former/former.reducers';

import { Former } from './Former';

export const FormerContainer: FC = () => {
  const dispatch = useAppDispatch();
  const history = useHistory();
  const { actions } = former;

  const {
    register,
    handleSubmit,
    setError,
    formState: { errors },
  } = useForm<FormData>();

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
      dispatch(actions.selectField({ fieldId }));
      dispatch(actions.resetFormErrors());
    },
    [dispatch],
  );

  const addField = useCallback(
    (kind: FieldKind) => {
      if (form) {
        dispatch(actions.addField({ formId: form.formId, kind }));
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
    (fieldId: string) => {
      dispatch(actions.resetFormErrors());
      dispatch(actions.selectField({ fieldId: undefined }));
    },
    [dispatch],
  );

  const saveForm = useCallback(async () => {
    if (ownerForm) {
      dispatch(actions.saveForm());
    } else if (formDefinition) {
      dispatch(actions.resetFormErrors());
      try {
        const data = await dispatch(postForm(formDefinition)).unwrap();
        history.push(`/browse/forms/${data.id}`);
      } catch (e: any) {
        dispatch(actions.setFormErrors({ errors: e?.details?.causes }));
      }
    }
  }, [dispatch, formDefinition, ownerForm]);

  if (!form) {
    return <></>;
  }

  return (
    <Former
      addField={addField}
      addOption={addOption}
      cancelField={(fieldId: string) => cancelField(fieldId)}
      errors={form.errors}
      fields={form.fields}
      formName={form.name}
      openSubForm={openSubForm}
      ownerFormName={ownerForm?.name}
      removeOption={removeOption}
      saveField={saveField}
      saveForm={saveForm}
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
  );
};

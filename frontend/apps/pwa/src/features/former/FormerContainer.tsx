import React, { FC, useCallback, useEffect } from 'react';
import { useHistory, useLocation } from 'react-router-dom';
import { FieldKind } from 'core-api-client';

import { useAppDispatch, useAppSelector } from '../../app/hooks';
import { fetchDatabases } from '../../reducers/database';
import { fetchFolders } from '../../reducers/folder';
import { fetchForms } from '../../reducers/form';
import {
  formerActions,
  formerGlobalSelectors,
  postForm,
} from '../../reducers/former';

import { Former } from './Former';

export const FormerContainer: FC = () => {
  const dispatch = useAppDispatch();
  const history = useHistory();

  // load data
  useEffect(() => {
    dispatch(formerActions.reset());
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
      dispatch(formerActions.setDatabase({ databaseId }));
    }
    const folderId = search.get('folderId');
    if (folderId) {
      dispatch(formerActions.setFolder({ folderId }));
    }
  }, [dispatch, location]);

  const setFormName = useCallback(
    (formName: string) => {
      if (form) {
        dispatch(formerActions.setFormName({ formId: form.formId, formName }));
      }
    },
    [dispatch, form],
  );

  const setSelectedField = useCallback(
    (fieldId: string | undefined) => {
      dispatch(formerActions.selectField({ fieldId }));
    },
    [dispatch],
  );

  const addField = useCallback(
    (kind: FieldKind) => {
      if (form) {
        dispatch(formerActions.addField({ formId: form.formId, kind }));
      }
    },
    [dispatch, form],
  );

  const setFieldOption = useCallback(
    (fieldId: string, i: number, value: string) => {
      dispatch(formerActions.setFieldOption({ fieldId, i, value }));
    },
    [dispatch],
  );

  const addOption = useCallback(
    (fieldId: string) => {
      dispatch(formerActions.addOption({ fieldId }));
    },
    [dispatch],
  );

  const removeOption = useCallback(
    (fieldId: string, i: number) => {
      dispatch(formerActions.removeOption({ fieldId, i }));
    },
    [dispatch],
  );

  const cancelField = useCallback(
    (fieldId: string) => {
      dispatch(formerActions.cancelFieldChanges({ fieldId }));
    },
    [dispatch],
  );

  const setFieldRequired = useCallback(
    (fieldId: string, required: boolean) => {
      dispatch(formerActions.setFieldRequired({ fieldId, required }));
    },
    [dispatch],
  );

  const setFieldIsKey = useCallback(
    (fieldId: string, isKey: boolean) => {
      dispatch(formerActions.setFieldIsKey({ fieldId, isKey }));
    },
    [dispatch],
  );

  const setFieldName = useCallback(
    (fieldId: string, name: string) => {
      dispatch(formerActions.setFieldName({ fieldId, name }));
    },
    [dispatch],
  );

  const setFieldDescription = useCallback(
    (fieldId: string, description: string) => {
      dispatch(formerActions.setFieldDescription({ fieldId, description }));
    },
    [dispatch],
  );

  const setFieldReferencedDatabaseId = useCallback(
    (fieldId: string, databaseId: string) => {
      dispatch(
        formerActions.setFieldReferencedDatabaseId({ fieldId, databaseId }),
      );
    },
    [dispatch],
  );

  const setFieldReferencedFormId = useCallback(
    (fieldId: string, formId: string) => {
      dispatch(formerActions.setFieldReferencedFormId({ fieldId, formId }));
    },
    [dispatch],
  );

  const openSubForm = useCallback(
    (fieldId: string) => {
      dispatch(formerActions.openSubForm({ fieldId }));
    },
    [dispatch],
  );

  const saveField = useCallback(
    (fieldId: string) => {
      dispatch(formerActions.selectField({ fieldId: undefined }));
    },
    [dispatch],
  );

  const saveForm = useCallback(async () => {
    if (ownerForm) {
      dispatch(formerActions.saveForm());
    } else if (formDefinition) {
      try {
        const data = await dispatch(postForm(formDefinition)).unwrap();
        history.push(`/browse/forms/${data.id}`);
        console.log('POST FORM', data);
      } catch (e) {
        console.log('POST FORM ERROR', e);
      }
    }
  }, [dispatch, formDefinition, ownerForm]);

  if (!form) {
    return <></>;
  }

  return (
    <Former
      formName={form.name}
      setFormName={setFormName}
      fields={form.fields}
      selectedFieldId={selectedField?.id}
      setSelectedField={setSelectedField}
      addField={addField}
      setFieldOption={setFieldOption}
      addOption={addOption}
      removeOption={removeOption}
      setFieldRequired={setFieldRequired}
      setFieldIsKey={setFieldIsKey}
      setFieldName={setFieldName}
      setFieldDescription={setFieldDescription}
      openSubForm={openSubForm}
      saveField={saveField}
      saveForm={saveForm}
      ownerFormName={ownerForm?.name}
      cancelField={(fieldId: string) => cancelField(fieldId)}
      setFieldReferencedDatabaseId={setFieldReferencedDatabaseId}
      setFieldReferencedFormId={setFieldReferencedFormId}
    />
  );
};

import React, { FC, useCallback, useEffect, useState } from 'react';
import { useHistory, useLocation, useParams } from 'react-router-dom';

import { useAppDispatch, useAppSelector } from '../../app/hooks';
import { fetchDatabases } from '../../reducers/database';
import { fetchFolders } from '../../reducers/folder';
import { fetchForms, selectRootForm } from '../../reducers/form';
import { recordActions } from '../../reducers/records';
import {
  postRecord,
  recorderActions,
  resetForm,
  selectCurrentForm,
  selectCurrentRecord,
  selectCurrentRootForm,
  selectPostRecords,
  selectSubRecords,
} from '../../reducers/recorder';

import { RecordEditorComponent } from './RecordEditor.component';

export const RecordEditorContainer: FC = () => {
  const dispatch = useAppDispatch();

  // load data
  useEffect(() => {
    dispatch(fetchDatabases());
    dispatch(fetchFolders());
    dispatch(fetchForms());
  }, [dispatch]);

  const params = useParams<{ formId: string }>();
  const location = useLocation();
  const history = useHistory();

  const [ownerRecordId, setOwnerRecordId] = useState<string | undefined>(
    undefined,
  );
  const formIdFromPath = params.formId;
  const currentRootForm = useAppSelector(selectCurrentRootForm);
  const rootFormFromPath = useAppSelector((s) =>
    selectRootForm(s, formIdFromPath),
  );
  const currentForm = useAppSelector(selectCurrentForm);
  const currentRecord = useAppSelector(selectCurrentRecord);
  const subRecords = useAppSelector((state) => {
    if (currentRecord) {
      return selectSubRecords(state, currentRecord.id);
    }
    return {};
  });

  useEffect(() => {
    const search = new URLSearchParams(location.search);
    const ownerRecordIdFromQryParam = search.get('ownerRecordId');
    if (ownerRecordIdFromQryParam !== ownerRecordId) {
      setOwnerRecordId(ownerRecordIdFromQryParam || undefined);
    }
  }, [ownerRecordId, location]);

  // make sure the form being edited is the one selected in the path
  useEffect(() => {
    if (!rootFormFromPath) {
      return;
    }

    if (rootFormFromPath.id !== currentRootForm?.id) {
      dispatch(
        resetForm({
          formId: formIdFromPath,
          ownerId: ownerRecordId,
        }),
      );
    }
  }, [
    dispatch,
    ownerRecordId,
    formIdFromPath,
    currentRootForm,
    rootFormFromPath,
  ]);

  const setFieldValue = useCallback(
    (key: string, value: any) => {
      if (currentRecord) {
        dispatch(
          recorderActions.setFieldValue({
            recordId: currentRecord.id,
            fieldId: key,
            value,
          }),
        );
      }
    },
    [dispatch, currentRecord],
  );

  const addSubRecord = useCallback(
    (ownerFieldId: string) => {
      if (!currentRecord) {
        return;
      }
      if (!currentForm) {
        return;
      }
      const field = currentForm.fields.find((f) => f.id === ownerFieldId);
      if (!field) {
        return;
      }
      if (!field.fieldType.subForm) {
        return;
      }
      const subFormId = field.id;

      dispatch(
        recorderActions.addSubRecord({
          formId: subFormId,
          ownerFieldId,
          ownerRecordId: currentRecord.id,
        }),
      );
    },
    [dispatch, currentForm, currentRecord],
  );

  const recordsToPost = useAppSelector(selectPostRecords);

  const saveRecord = useCallback(async () => {
    // do not save if we are not positioned on a record (should not happen)
    if (!currentRecord) {
      return;
    }
    // do not save if we are not positioned on a form (should not happen)
    if (!currentForm) {
      return;
    }

    if (currentRecord.formId !== formIdFromPath) {
      if (currentRecord.ownerId) {
        dispatch(
          recorderActions.selectRecord({
            recordId: currentRecord.ownerId,
          }),
        );
      }
    } else {
      const recordResponse = await dispatch(postRecord(recordsToPost)).unwrap();
      dispatch(recordActions.addMany(recordResponse));
      history.push(`/browse/records/${recordResponse[0].id}`);
    }
  }, [dispatch, formIdFromPath, currentRecord, recordsToPost, currentForm]);

  const selectSubRecord = useCallback(
    (subRecordId: string) => {
      dispatch(recorderActions.selectRecord({ recordId: subRecordId }));
    },
    [dispatch],
  );

  if (!currentForm) {
    return <></>;
  }

  if (!currentRecord) {
    return <></>;
  }

  return (
    <RecordEditorComponent
      setValue={setFieldValue}
      fields={currentForm?.fields}
      values={currentRecord?.values}
      addSubRecord={addSubRecord}
      saveRecord={saveRecord}
      subRecords={subRecords}
      selectSubRecord={selectSubRecord}
    />
  );
};

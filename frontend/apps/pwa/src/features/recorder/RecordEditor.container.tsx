import React, { FC, useCallback, useEffect, useState } from 'react';
import { useHistory, useLocation, useParams } from 'react-router-dom';
import { useForm, FormProvider } from 'react-hook-form';

import { useAppDispatch, useAppSelector } from '../../app/hooks';
import { fetchDatabases } from '../../reducers/database';
import { fetchFolders } from '../../reducers/folder';
import { fetchForms, selectRootForm } from '../../reducers/form';
import { recordActions } from '../../reducers/records';
import {
  postRecord,
  resetForm,
} from '../../reducers/Recorder/recorder.reducers';
import recorderSlice from '../../reducers/Recorder';
import {
  selectCurrentForm,
  selectCurrentRecord,
  selectCurrentRootForm,
  selectPostRecords,
  selectSubRecords,
} from '../../reducers/Recorder/recorder.selectors';
import { FormValue } from '../../reducers/Recorder/types';

import { RecordEditorComponent } from './RecordEditor.component';

export const RecordEditorContainer: FC = () => {
  const dispatch = useAppDispatch();

  const { actions } = recorderSlice;

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

  const useFormObject = useForm<FormValue>({
    defaultValues: currentRecord,
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

  const handleFieldValueChange = useCallback(
    (key: string, value: any) => {
      if (currentRecord) {
        dispatch(
          actions.setFieldValue({
            recordId: currentRecord.id,
            fieldId: key,
            value,
          }),
        );
      }
    },
    [dispatch, currentRecord],
  );

  const handleAddSubRecord = useCallback(
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
        actions.addSubRecord({
          formId: subFormId,
          ownerFieldId,
          ownerRecordId: currentRecord.id,
        }),
      );
    },
    [dispatch, currentForm, currentRecord],
  );

  const recordsToPost = useAppSelector(selectPostRecords);

  const handleSaveRecord = useCallback(async () => {
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
          actions.selectRecord({
            recordId: currentRecord.ownerId,
          }),
        );
      }
    } else {
      try {
        const recordResponse = await dispatch(
          postRecord(recordsToPost),
        ).unwrap();
        dispatch(recordActions.addMany(recordResponse));
        history.push(`/browse/records/${recordResponse[0].id}`);
      } catch (e: any) {}
    }
  }, [dispatch, formIdFromPath, currentRecord, recordsToPost, currentForm]);

  const handleSelectSubRecord = useCallback(
    (subRecordId: string) => {
      dispatch(actions.selectRecord({ recordId: subRecordId }));
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
    <FormProvider {...useFormObject}>
      <RecordEditorComponent
        onChangeValue={handleFieldValueChange}
        fields={currentForm?.fields}
        values={currentRecord?.values}
        errors={useFormObject.formState.errors}
        onAddSubRecord={handleAddSubRecord}
        onSaveRecord={useFormObject.handleSubmit(handleSaveRecord)}
        subRecords={subRecords}
        onSelectSubRecord={handleSelectSubRecord}
      />
    </FormProvider>
  );
};

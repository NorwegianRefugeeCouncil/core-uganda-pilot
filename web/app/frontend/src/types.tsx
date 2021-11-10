import { StackScreenProps } from '@react-navigation/stack';
import React from 'react';

import { RecordsAction, RecordsStoreProps } from './reducers/recordsReducers';

export type StackParamList = {
    forms: undefined;
    records: { formId: string; databaseId: string };
    addRecord: { formId: string; recordId: string };
    viewRecord: { formId: string; recordId: string };
    designSystem: undefined;
};

export type FormsScreenProps = StackScreenProps<StackParamList, 'forms'>;
export type RecordsScreenProps = StackScreenProps<StackParamList, 'records'> & {
    state: RecordsStoreProps;
    dispatch: React.Dispatch<RecordsAction>;
};

import { StackScreenProps } from "@react-navigation/stack";
import React from "react";

import { RecordsAction, RecordsStoreProps } from "../reducers/recordsReducers";

export type StackParamList = {
    forms: undefined;
    records: { formId: string; databaseId: string };
    addRecord: { formId: string; recordId: string };
    viewRecord: { formId: string; recordId: string };
    designSystem: undefined;
};

type ReducerProps = {
    state: RecordsStoreProps;
    dispatch: React.Dispatch<RecordsAction>;
}

export type FormsScreenContainerProps = StackScreenProps<StackParamList, 'forms'>;
export type RecordsScreenContainerProps = StackScreenProps<StackParamList, 'records'> & ReducerProps;
export type AddRecordScreenContainerProps = StackScreenProps<StackParamList, 'addRecord'> & ReducerProps;
export type ViewRecordScreenContainerProps = StackScreenProps<StackParamList, 'viewRecord'> & ReducerProps
export type DesignSystemScreenProps = StackScreenProps<StackParamList, 'designSystem'>

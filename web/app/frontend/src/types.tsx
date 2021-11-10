import { StackScreenProps } from '@react-navigation/stack';

export type StackParamList = {
    forms: undefined;
    records: { formId: string; databaseId: string } | undefined;
    addRecord: { formId: string; recordId: string };
    viewRecord: { formId: string; recordId: string };
    designSystem: undefined;
};

export type FormsScreenProps = StackScreenProps<StackParamList, 'forms'>;

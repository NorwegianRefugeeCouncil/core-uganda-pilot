export type FieldValue = {
  fieldId: string;
  value: string | string[] | null;
};

export type Record = {
  id: string;
  databaseId: string;
  formId: string;
  ownerId: string | undefined;
  values: FieldValue[];
};

export type LocalRecord = Record & {
  isNew: boolean;
};

export type RecordList = { items: Record[] };

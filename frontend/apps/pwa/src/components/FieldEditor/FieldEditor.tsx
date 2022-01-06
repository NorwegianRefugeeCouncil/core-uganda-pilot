import React, { FC } from 'react';

import { FieldEditorProps } from '../formFields/types';
import {
  TextFieldEditor,
  DateFieldEditor,
  MonthFieldEditor,
  MultilineTextFieldEditor,
  MultiSelectFieldEditor,
  QuantityFieldEditor,
  WeekFieldEditor,
  SingleSelectFieldEditor,
  SubFormFieldEditor,
  ReferenceFieldEditor,
} from '../formFields';

export const FieldEditor: FC<FieldEditorProps> = (props) => {
  const {
    field: { fieldType },
  } = props;
  if (fieldType.text) {
    return <TextFieldEditor {...props} />;
  }
  if (fieldType.week) {
    return <WeekFieldEditor {...props} />;
  }
  if (fieldType.subForm) {
    return <SubFormFieldEditor {...props} />;
  }
  if (fieldType.reference) {
    return <ReferenceFieldEditor {...props} />;
  }
  if (fieldType.multilineText) {
    return <MultilineTextFieldEditor {...props} />;
  }
  if (fieldType.date) {
    return <DateFieldEditor {...props} />;
  }
  if (fieldType.month) {
    return <MonthFieldEditor {...props} />;
  }
  if (fieldType.quantity) {
    return <QuantityFieldEditor {...props} />;
  }
  if (fieldType.singleSelect) {
    return <SingleSelectFieldEditor {...props} />;
  }
  if (fieldType.multiSelect) {
    return <MultiSelectFieldEditor {...props} />;
  }
  return <></>;
};

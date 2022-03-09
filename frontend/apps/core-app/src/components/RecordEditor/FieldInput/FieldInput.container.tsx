import * as React from 'react';
import {
  FormDefinition,
  Record,
  getFieldKind,
  FieldKind,
} from 'core-api-client';

import { TextFieldInput } from './TextFieldInput.component';
import { MultilineTextFieldInput } from './MultilineTextFieldInput.component';
import { ReferenceFieldInput } from './ReferenceFieldInput.component';
import { DateFieldInput } from './DateFieldInput.component';
import { MonthFieldInput } from './MonthFieldInput.component';
import { WeekFieldInput } from './WeekFieldInput.component';
import { SingleSelectFieldInput } from './SingleSelectFieldInput.component';
import { CheckboxFieldInput } from './CheckboxFieldInput.component';
import { SubFormFieldInput } from './SubFormFieldInput.component';
import { MultiSelectFieldInput } from './MultiSelectFieldInput.component';

type Props = {
  field: FormDefinition['fields'][number];
};

export const FieldInput: React.FC<Props> = ({ field }) => {
  switch (getFieldKind(field.fieldType)) {
    case FieldKind.Text:
      return <TextFieldInput />;
    case FieldKind.MultilineText:
      return <MultilineTextFieldInput />;
    case FieldKind.Reference:
      return <ReferenceFieldInput />;
    case FieldKind.Date:
      return <DateFieldInput />;
    case FieldKind.Month:
      return <MonthFieldInput />;
    case FieldKind.Week:
      return <WeekFieldInput />;
    case FieldKind.SingleSelect:
      return <SingleSelectFieldInput />;
    case FieldKind.MultiSelect:
      return <MultiSelectFieldInput />;
    case FieldKind.Checkbox:
      return <CheckboxFieldInput />;
    case FieldKind.SubForm:
      return <SubFormFieldInput />;
    default:
      return null;
  }
};

import * as React from 'react';
import {
  FormDefinition,
  Record,
  getFieldKind,
  FieldKind,
} from 'core-api-client';

import { TextFieldInput } from './TextFieldInput.component';
import { MultilineTextFieldInput } from './MultilineTextFieldInput.component';
import { QuantityFieldInput } from './QuantityFieldInput.component';
import { ReferenceFieldInput } from './ReferenceFieldInput.component';
import { DateFieldInput } from './DateFieldInput.component';
import { MonthFieldInput } from './MonthFieldInput.component';
import { WeekFieldInput } from './WeekFieldInput.component';
import { SingleSelectFieldInput } from './SingleSelectFieldInput.component';
import { CheckboxFieldInput } from './CheckboxFieldInput.component';
import { SubFormFieldInput } from './SubFormFieldInput.component';
import { MultiSelectFieldInput } from './MultiSelectFieldInput.component';

type Props = {
  form: FormDefinition;
  field: FormDefinition['fields'][number];
};

export const FieldInput: React.FC<Props> = ({ form, field }) => {
  switch (getFieldKind(field.fieldType)) {
    case FieldKind.Text:
      return <TextFieldInput formId={form.id} field={field} />;
    case FieldKind.MultilineText:
      return <MultilineTextFieldInput formId={form.id} field={field} />;
    case FieldKind.Quantity:
      return <QuantityFieldInput formId={form.id} field={field} />;
    case FieldKind.Reference:
      return (
        <ReferenceFieldInput
          formId={form.id}
          databaseId={form.databaseId}
          field={field}
        />
      );
    case FieldKind.Date:
      return <DateFieldInput formId={form.id} field={field} />;
    case FieldKind.Month:
      return <MonthFieldInput formId={form.id} field={field} />;
    case FieldKind.Week:
      return <WeekFieldInput formId={form.id} field={field} />;
    case FieldKind.SingleSelect:
      return <SingleSelectFieldInput formId={form.id} field={field} />;
    case FieldKind.MultiSelect:
      return <MultiSelectFieldInput formId={form.id} field={field} />;
    case FieldKind.Checkbox:
      return <CheckboxFieldInput formId={form.id} field={field} />;
    case FieldKind.SubForm:
      return <SubFormFieldInput form={form} field={field} />;
    default:
      return null;
  }
};

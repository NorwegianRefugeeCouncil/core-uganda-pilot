import * as React from 'react';
import { FormGroup } from './form-group.component';
import { FormLabel } from './form-label.component';
import { FormControl } from './form-control.component';
import { FormSelect } from './form-select.component';
import { FormText } from './form-text.component';
import { classNames, Box, PolymorphicComponentProps } from '@ui-helpers/utils';
import { FormCheck } from './form-check.component';

export type FormOwnProps = {
  inline?: true;
  validated?: true;
};

export type FormProps<E extends React.ElementType> = PolymorphicComponentProps<
  E,
  FormOwnProps
>;

export type FormStatic = {
  Group?: typeof FormGroup;
  Label?: typeof FormLabel;
  Control?: typeof FormControl;
  Select?: typeof FormSelect;
  Check?: typeof FormCheck;
  Text?: typeof FormText;
};

type Form = <E extends React.ElementType = 'form'>(
  props: FormProps<E>
) => React.ReactElement | null;

export const Form: Form & FormStatic = React.forwardRef(
  <E extends React.ElementType = 'form'>(
    { as, inline, validated, className: customClass, ...rest }: FormProps<E>,
    ref: typeof rest.ref
  ) => {
    const className = classNames(customClass, {
      'row row-cols-lg-auto g-3 align-items-center': inline,
      'was-validated': validated,
      'needs-validation': !validated,
    });
    return <Box ref={ref} className={className} {...rest} />;
  }
);

Form.Group = FormGroup;
Form.Label = FormLabel;
Form.Control = FormControl;
Form.Select = FormSelect;
Form.Check = FormCheck;
Form.Text = FormText;

import * as React from 'react';
import { FormGroup } from './form-group.component';
import { FormLabel } from './form-label.component';
import { FormControl } from './form-control.component';
import { FormSelect } from './form-select.component';
import { FormText } from './form-text.component';
import { FormCheck } from './form-check.component';
import { classNames } from '@core/ui-toolkit/util/utils';

export interface FormProps extends React.ComponentPropsWithRef<'form'> {
  inline?: true;
  validated?: true;
}

export type FormStatic = {
  Group?: typeof FormGroup;
  Label?: typeof FormLabel;
  Control?: typeof FormControl;
  Select?: typeof FormSelect;
  Check?: typeof FormCheck;
  Text?: typeof FormText;
};

export type Form = React.ForwardRefExoticComponent<
  React.PropsWithRef<FormProps>
> &
  FormStatic;

export const Form: Form = React.forwardRef(
  ({ inline, validated, className: customClass, ...rest }, ref) => {
    const className = classNames(customClass, {
      'row row-cols-lg-auto g-3 align-items-center': inline,
      'was-validated': validated,
      'needs-validation': !validated,
    });
    return <form ref={ref} className={className} {...rest} />;
  }
);

Form.Group = FormGroup;
Form.Label = FormLabel;
Form.Control = FormControl;
Form.Select = FormSelect;
Form.Check = FormCheck;
Form.Text = FormText;

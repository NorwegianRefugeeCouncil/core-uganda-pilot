import * as React from 'react';
import FormGroup from './form-group.component';
import FormLabel from './form-label.component';
import FormControl from './form-control.component';
import FormSelect from './form-select.component';
import { classNames } from '@ui-helpers/utils';

interface Props<C extends React.ElementType> {
  as?: C;
  children: React.ReactNode;
  inline?: true;
  validated?: boolean;
}

type FormProps<C extends React.ElementType> = Props<C> &
  Omit<React.ComponentPropsWithRef<C>, keyof Props<C>> & {
    Group: typeof FormGroup;
    Label: typeof FormLabel;
    Control: typeof FormControl;
    Select: typeof FormSelect;
  };

const Form = <C extends React.ElementType = 'form'>({
  as,
  inline,
  validated,
  className: customClass,
  ...rest
}) => {
  const className = classNames(customClass, {
    'row row-cols-lg-auto g-3 align-items-center': inline,
    'was-validated': validated,
    'needs-validation': !validated,
  });
  const Component = as ?? 'form';
  return <Component className={className} {...rest} />;
};

Form.displayName = 'Form';

Form.Group = FormGroup;
Form.Label = FormLabel;
Form.Control = FormControl;
Form.Select = FormSelect;

export default Form;

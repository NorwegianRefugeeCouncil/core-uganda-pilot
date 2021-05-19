import { classNames } from '@ui-helpers/utils';
import * as React from 'react';
import { FormContext } from './form-context';

interface FormLabelProps extends React.ComponentPropsWithoutRef<'label'> {
  srOnly?: true;
  htmlFor?: string;
}

export const FormLabel: React.FC<FormLabelProps> = ({
  srOnly,
  htmlFor,
  className: customClass,
  ...rest
}) => {
  const { controlId } = React.useContext(FormContext);
  const className = classNames(customClass, {
    'visually-hidden': srOnly,
  });
  return (
    <label htmlFor={htmlFor ?? controlId} className={className} {...rest} />
  );
};

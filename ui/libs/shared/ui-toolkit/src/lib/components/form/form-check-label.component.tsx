import { classNames } from '@ui-helpers/utils';
import * as React from 'react';
import FormContext from './form-context';

interface FormCheckLabelProps extends React.ComponentPropsWithoutRef<'label'> {
  srOnly?: true;
  htmlFor?: string;
}

const FormCheckLabel: React.FC<FormCheckLabelProps> = ({
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

export default FormCheckLabel;

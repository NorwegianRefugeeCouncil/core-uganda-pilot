import * as React from 'react';
import { FormContext } from './form-context';
import { classNames } from '@core/ui-toolkit/util/utils';

interface FormCheckLabelProps extends React.ComponentPropsWithoutRef<'label'> {
  srOnly?: true;
  htmlFor?: string;
}

export const FormCheckLabel: React.FC<FormCheckLabelProps> = ({
  srOnly,
  htmlFor,
  className: customClass,
  ...rest
}) => {
  const { controlId } = React.useContext(FormContext);
  const className = classNames(
    'form-check-label',
    {
      'visually-hidden': srOnly,
    },
    customClass
  );
  return (
    <label htmlFor={htmlFor ?? controlId} className={className} {...rest} />
  );
};

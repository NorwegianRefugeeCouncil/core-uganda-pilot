import { classNames } from '@ui-helpers/utils';
import * as React from 'react';
import FormContext from './form-context';

interface Props<C extends React.ElementType> {
  type: 'text' | 'textarea' | 'email' | 'password' | 'file' | 'color';
  displaySize?: 'sm' | 'lg';
  plaintext?: true;
  isValid?: true;
  isInvalid?: true;
  children: React.ReactNode;
}

type FormControlProps<C extends React.ElementType> = Props<C> &
  Omit<React.ComponentPropsWithRef<C>, keyof Props<C>>;

const FormControl = <C extends React.ElementType = 'input'>({
  id,
  displaySize,
  type,
  plaintext,
  isValid,
  isInvalid,
  className: customClass,
  ...rest
}) => {
  const Component = type === 'textarea' ? 'textarea' : type;
  const { controlId } = React.useContext(FormContext);
  const className = classNames(
    'form-control',
    {
      [`form-control-${displaySize}`]: displaySize != null,
      'form-control-plaintext': rest.readOnly && plaintext,
      'is-valid': isValid,
      'is-invalid': isInvalid,
    },
    customClass
  );
  return <Component id={id ?? controlId} className={className} {...rest} />;
};
export default FormControl;

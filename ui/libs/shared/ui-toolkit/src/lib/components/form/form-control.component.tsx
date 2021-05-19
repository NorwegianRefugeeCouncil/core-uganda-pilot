import * as React from 'react';
import { FormContext } from './form-context';
import { classNames } from '@ui-helpers/utils';

export interface BaseFormControlProps {
  type: 'text' | 'textarea' | 'email' | 'password' | 'file' | 'color';
  displaySize?: 'sm' | 'lg';
  plaintext?: true;
  isValid?: true;
  isInvalid?: true;
}

type InputProps = BaseFormControlProps & JSX.IntrinsicElements['input'];
type TextareaProps = BaseFormControlProps & JSX.IntrinsicElements['textarea'];

export type FormControl = {
  (props: InputProps): JSX.Element;
  (props: TextareaProps): JSX.Element;
};

function isPropsForTextarea(
  props: InputProps | TextareaProps
): props is TextareaProps {
  return 'cols' in props;
}

export const FormControl: FormControl = (props: InputProps | TextareaProps) => {
  const {
    id,
    displaySize,
    type,
    plaintext,
    isValid,
    isInvalid,
    className: customClass,
    ...rest
  } = props;
  const { controlId } = React.useContext(FormContext);
  const className = classNames(
    {
      'form-control': !plaintext,
      'form-control-plaintext': rest.readOnly && plaintext,
      [`form-control-${displaySize}`]: displaySize != null,
      [`form-control-color`]: type === 'color',
      'is-valid': isValid,
      'is-invalid': isInvalid,
    },
    customClass
  );
  if (type === 'textarea')
    return (
      <textarea
        id={id ?? controlId}
        className={className}
        {...(rest as TextareaProps)}
      />
    );
  else
    return (
      <input
        type={type}
        id={id ?? controlId}
        className={className}
        {...(rest as InputProps)}
      />
    );
};

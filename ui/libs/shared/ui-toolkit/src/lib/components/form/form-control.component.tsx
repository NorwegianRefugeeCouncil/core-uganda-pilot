import * as React from 'react';
import { FormContext } from './form-context';
import { classNames } from '@ui-helpers/utils';
import { BsInputTypes, NonBsInputTypes } from '@ui-helpers/types';

export interface BaseFormControlProps {
  type?: BsInputTypes | NonBsInputTypes;
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
    type = 'text',
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
  else if (type === 'number' || type === 'date')
    return (
      <input
        type={type}
        id={id ?? controlId}
        className={customClass}
        {...(rest as InputProps)}
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

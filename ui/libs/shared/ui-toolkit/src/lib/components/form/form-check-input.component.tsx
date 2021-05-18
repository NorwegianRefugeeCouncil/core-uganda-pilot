import * as React from 'react';
import classNames from 'classnames';
import FormContext from './form-context';

interface FormCheckInputProps extends React.ComponentPropsWithRef<'input'> {
  type: 'radio' | 'checkbox';
  label?: string;
  inline?: true;
  isValid?: true;
  isInvalid?: true;
}

const FormCheckInput = React.forwardRef<HTMLInputElement, FormCheckInputProps>(
  ({ id, type, isValid, isInvalid, className: customClass, ...rest }, ref) => {
    const { controlId } = React.useContext(FormContext);
    const className = classNames('form-check-input', customClass);
    return (
      <input ref={ref} id={id ?? controlId} className={className} {...rest} />
    );
  }
);

export default FormCheckInput;

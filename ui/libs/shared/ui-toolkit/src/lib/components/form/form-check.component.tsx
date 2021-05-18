import * as React from 'react';
import FormContext, { FormContextInterface } from './form-context';
import FormCheckLabel from './form-check-label.component';
import FormCheckInput from './form-check-input.component';
import { classNames } from '@ui-helpers/utils';

export interface FormCheckProps extends React.ComponentPropsWithRef<'div'> {
  id: string;
  label?: string;
  type?: 'radio' | 'checkbox' | 'switch';
  inline?: true;
  isValid?: true;
  isInvalid?: true;
  disabled?: true;
  required?: true;
}

interface FormCheckStatic {
  Label: typeof FormCheckLabel;
  Input: typeof FormCheckInput;
}

type FormCheckComponent = React.ForwardRefExoticComponent<
  React.PropsWithRef<FormCheckProps>
> &
  Partial<FormCheckStatic>;

const FormCheck = React.forwardRef<HTMLInputElement, FormCheckProps>(
  (
    {
      id,
      label,
      type = 'checkbox',
      inline,
      isValid,
      isInvalid,
      disabled,
      required,
      className: customClass,
      children,
      ...rest
    },
    ref
  ) => {
    const { controlId } = React.useContext(FormContext);
    const innerFormContext: FormContextInterface = {
      controlId: id + controlId,
    };
    const className = classNames(
      'form-check',
      {
        'form-check-inline': inline,
        'form-check-switch': type === 'switch',
        'is-invalid': isInvalid,
      },
      customClass
    );

    const labelComponent = label ? (
      <FormCheckLabel>{label}</FormCheckLabel>
    ) : null;
    const inputComponent = (
      <FormCheckInput
        ref={ref}
        type={type === 'switch' ? 'checkbox' : type}
        disabled={disabled}
        required={required}
        isValid={isValid}
        isInvalid={isInvalid}
      />
    );

    return (
      <FormContext.Provider value={innerFormContext}>
        <div ref={ref} id={id ?? controlId} className={className} {...rest}>
          {children ?? (
            <>
              {labelComponent}
              {inputComponent}
            </>
          )}
        </div>
      </FormContext.Provider>
    );
  }
) as FormCheckComponent;

FormCheck.displayName = 'FormCheck';

FormCheck.Label = FormCheckLabel;
FormCheck.Input = FormCheckInput;

export default FormCheck;

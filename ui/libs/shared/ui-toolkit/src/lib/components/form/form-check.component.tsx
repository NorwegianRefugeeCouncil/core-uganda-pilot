import * as React from 'react';
import { FormContext, FormContextInterface } from './form-context';
import { FormCheckLabel } from './form-check-label.component';
import { FormCheckInput } from './form-check-input.component';
import { classNames } from '@core/shared/ui-toolkit/util/utils';

export interface FormCheckProps extends React.ComponentPropsWithRef<'div'> {
  id: string;
  label?: string;
  type?: 'radio' | 'checkbox' | 'switch';
  name?: string;
  defaultChecked?: true;
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
      name,
      defaultChecked,
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
    // const { controlId } = React.useContext(FormContext);
    const innerFormContext: FormContextInterface = {
      controlId: id,
    };

    const className = classNames(
      'form-check',
      {
        'form-check-inline': inline,
        'form-switch': type === 'switch',
        'is-invalid': isInvalid,
      },
      customClass
    );

    const inputComponent = (
      <FormCheckInput
        id={id}
        type={type === 'switch' ? 'checkbox' : type}
        name={name}
        disabled={disabled}
        required={required}
        isValid={isValid}
        isInvalid={isInvalid}
        defaultChecked={defaultChecked}
      />
    );

    const labelComponent = label ? (
      <FormCheckLabel htmlFor={id}>{label}</FormCheckLabel>
    ) : null;

    return (
      <FormContext.Provider value={innerFormContext}>
        <div ref={ref} className={className} {...rest}>
          {children ?? (
            <>
              {inputComponent}
              {labelComponent}
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

export { FormCheck };

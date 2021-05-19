import * as React from 'react';
import { classNames } from '@ui-helpers/utils';

export interface FormSelectProps extends React.ComponentPropsWithRef<'select'> {
  options?: Array<{ value: any; label?: string; disabled?: true }>;
  selectedOptionIdx?: number;
  displaySize?: 'sm' | 'lg';
}

export const FormSelect = React.forwardRef<HTMLSelectElement, FormSelectProps>(
  (
    {
      options,
      selectedOptionIdx,
      displaySize,
      className: customClass,
      children,
      ...rest
    },
    ref
  ) => {
    const className = classNames(
      'form-select',
      {
        [`form-select-${displaySize}`]: displaySize != null,
      },
      customClass
    );
    return (
      <select className={className} {...rest}>
        {children ??
          options.map((option, idx) => (
            <option
              value={option.value}
              selected={idx === selectedOptionIdx}
              disabled={option.disabled}
            >
              {option.label ?? option.value ?? `option ${idx}`}
            </option>
          ))}
      </select>
    );
  }
);

FormSelect.displayName = 'FormSelect';

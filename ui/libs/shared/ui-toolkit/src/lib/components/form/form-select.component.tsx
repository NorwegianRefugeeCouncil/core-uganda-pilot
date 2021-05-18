import * as React from 'react';
import { classNames } from '@ui-helpers/utils';
import { optionCSS } from 'react-select/src/components/Option';

export interface FormSelectProps extends React.ComponentPropsWithRef<'select'> {
  options?: Array<{ value: any; label?: string; disabled?: true }>;
  selectedOptionIdx?: number;
  displaySize?: 'sm' | 'lg';
}

const FormSelect = React.forwardRef<HTMLSelectElement, FormSelectProps>(
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
      <select {...rest}>
        {children ??
          options.map((option, idx) => (
            <option value={option.value} selected={idx === selectedOptionIdx}>
              {option.label ?? option.value ?? `option ${idx}`}
            </option>
          ))}
      </select>
    );
  }
);

export default FormSelect;

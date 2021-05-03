import {
  forwardRef,
  FunctionComponent,
  InputHTMLAttributes,
  LabelHTMLAttributes,
  SelectHTMLAttributes,
} from 'react';
import classNames from 'classnames';
import { CollapseProps } from '../collapse/collapse.component';

type SelectProps = SelectHTMLAttributes<HTMLSelectElement>;

export const FormSelect = forwardRef<HTMLSelectElement, SelectProps>(
  (props, ref) => {
    return (
      <select
        {...props}
        ref={ref}
        className={classNames(props.className, 'form-select')}
      />
    );
  }
);

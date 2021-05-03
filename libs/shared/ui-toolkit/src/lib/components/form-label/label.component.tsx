import { forwardRef, LabelHTMLAttributes } from 'react';
import classNames from 'classnames';

export type FormLabelProps = LabelHTMLAttributes<HTMLLabelElement>

export const FormLabel = forwardRef<HTMLLabelElement, FormLabelProps>(
  (props, ref) => {
    const { className, ...otherProps } = props;
    return (<label {...otherProps} ref={ref} className={classNames(className, 'form-label')}>{props.children}</label>);
  });

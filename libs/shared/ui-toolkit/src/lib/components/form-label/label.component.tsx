import { FunctionComponent, LabelHTMLAttributes } from 'react';
import classNames from 'classnames';

export type FormLabelProps = LabelHTMLAttributes<HTMLLabelElement>;

export const FormLabel: FunctionComponent<FormLabelProps> = (props) => {
  const { className, ...otherProps } = props;
  return (
    <label {...otherProps} className={classNames(className, 'form-label')}>
      {props.children}
    </label>
  );
};

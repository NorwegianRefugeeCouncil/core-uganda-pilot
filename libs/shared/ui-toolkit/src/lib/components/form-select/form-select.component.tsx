import { FunctionComponent, SelectHTMLAttributes } from 'react';
import classNames from 'classnames';

type SelectProps = SelectHTMLAttributes<HTMLSelectElement>;

export const FormSelect: FunctionComponent<SelectProps> = (props) => {
  return (
    <select {...props} className={classNames(props.className, 'form-select')} />
  );
};

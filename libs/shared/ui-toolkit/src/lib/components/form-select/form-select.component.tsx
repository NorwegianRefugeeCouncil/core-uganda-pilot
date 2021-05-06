import * as React from 'react'
import classNames from 'classnames';

type SelectProps = React.ComponentPropsWithRef<'select'>;

export const FormSelect: React.FC<SelectProps> = (props) => {
  return (
    <select {...props} className={classNames(props.className, 'form-select')} />
  );
};

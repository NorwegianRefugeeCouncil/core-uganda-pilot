import * as React from 'react';
import classNames from 'classnames';

export type FormLabelProps = React.ComponentPropsWithRef<'label'>;

export const FormLabel: React.FC<FormLabelProps> = (props) => {
  const { className, ...otherProps } = props;
  return (
    <label {...otherProps} className={classNames(className, 'form-label')}>
      {props.children}
    </label>
  );
};

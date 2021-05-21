import * as React from 'react';
import { classNames } from '@core/ui-toolkit/util/utils';

export interface ButtonGroupProps
  extends React.ComponentPropsWithoutRef<'div'> {
  size?: 'sm' | 'lg';
}

type ButtonGroup = React.FC<ButtonGroupProps>;

export const ButtonGroup: ButtonGroup = ({
  size,
  className: customClass,
  children,
  ...rest
}) => {
  const className = classNames('btn-group', customClass, {
    [`btn-group-${size}`]: size != null,
  });
  return (
    <div className={className} role="group" {...rest}>
      {children}
    </div>
  );
};

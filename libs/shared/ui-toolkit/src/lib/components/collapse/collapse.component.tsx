import * as React from 'react';
import classNames from 'classnames';

export type CollapseProps = React.ComponentPropsWithRef<'div'> & {
  show?: boolean;
};

export const Collapse: React.FC<CollapseProps> = ({
  show = true,
  className: customClass,
  children,
  ...rest
}) => {
  const classeName = classNames(customClass, 'collapse', { show });
  return (
    <div {...rest} className={classeName}>
      {children}
    </div>
  );
};

/* eslint-disable @typescript-eslint/no-empty-interface */
import * as React from 'react';
import { classNames } from '@core/shared/ui-toolkit/util/utils';

export interface NavItemProps extends React.ComponentPropsWithoutRef<'li'> {
  dropdown?: boolean;
}

export const NavItem: React.FC<NavItemProps> = ({
  dropdown = false,
  className: customClass,
  children,
  ...rest
}) => {
  const className = classNames('nav-item', customClass, { dropdown });
  return (
    <li className={className} {...rest}>
      {children}
    </li>
  );
};

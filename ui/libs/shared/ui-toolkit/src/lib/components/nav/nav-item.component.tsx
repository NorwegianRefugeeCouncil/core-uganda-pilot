/* eslint-disable @typescript-eslint/no-empty-interface */
import * as React from 'react';
import { classNames } from '../../helpers/utils';

export interface NavItemProps extends React.ComponentPropsWithoutRef<'li'> {
  dropdown?: boolean;
}

const NavItem: React.FC<NavItemProps> = ({
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

export default NavItem;

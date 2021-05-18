/* eslint-disable @typescript-eslint/no-empty-interface */
import * as React from 'react';
import { classNames } from '../../helpers/utils';

export interface NavLinkProps extends React.ComponentPropsWithoutRef<'a'> {
  isActive?: boolean;
  isDisabled?: boolean;
}

const NavLink: React.FC<NavLinkProps> = ({
  isActive: active = false,
  isDisabled: disabled = false,
  className: customClass,
  children,
  ...rest
}) => {
  const className = classNames('nav-link', customClass, { active, disabled });
  return (
    <a className={className} {...rest}>
      {children}
    </a>
  );
};

export default NavLink;

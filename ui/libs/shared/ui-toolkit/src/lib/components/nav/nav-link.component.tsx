/* eslint-disable @typescript-eslint/no-empty-interface */
import * as React from 'react';
import { classNames } from '@core/ui-toolkit/util/utils';

export interface NavLinkProps extends React.ComponentPropsWithoutRef<'a'> {
  active?: boolean;
  disabled?: boolean;
}

export const NavLink: React.FC<NavLinkProps> = ({
  active = false,
  disabled = false,
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

/* eslint-disable @typescript-eslint/no-empty-interface */
import * as React from 'react';
import { classNames } from '@core/ui-toolkit/util/utils';
import { Size, Theme } from '@core/ui-toolkit/util/types';

export interface NavProps extends React.ComponentPropsWithoutRef<'ul'> {}

const Nav: React.FC<NavProps> = ({
  className: customClass,
  children,
  ...rest
}) => {
  const className = classNames('navbar-nav', customClass);
  return (
    <ul className={className} {...rest}>
      {children}
    </ul>
  );
};

export interface BrandProps extends React.ComponentPropsWithoutRef<'a'> {}
const Brand: React.FC<BrandProps> = ({
  className: customClass,
  children,
  ...rest
}) => {
  const className = classNames('navbar-brand', customClass);
  return (
    <a className={className} {...rest}>
      {children}
    </a>
  );
};

export interface NavbarProps extends React.ComponentPropsWithoutRef<'div'> {
  canCollapse?: boolean;
  theme?: Theme;
  size?: Size;
}

export const Navbar: React.FC<NavbarProps> = ({
  canCollapse = true,
  size,
  theme = 'light',
  className: customClass,
  children,
  ...rest
}) => {
  const className = classNames(
    'navbar',
    `navbar-${theme}`,
    `bg-${theme}`,
    customClass,
    {
      [`navbar-expand-${size}`]: canCollapse,
    }
  );
  return (
    <nav className={className} {...rest}>
      <div className="container-fluid">{children}</div>
    </nav>
  );
};

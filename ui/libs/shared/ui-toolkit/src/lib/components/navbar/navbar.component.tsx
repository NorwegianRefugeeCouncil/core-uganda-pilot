/* eslint-disable @typescript-eslint/no-empty-interface */
import * as React from 'react';
import { classNames } from '@core/ui-toolkit/util/utils';
import { Color, Size, Theme } from '@core/ui-toolkit/util/types';
import { NavbarBrand } from './navbar-brand.component';
import { NavbarNav } from './navbar-nav.component';
import { NavbarText } from './navbar-text.component';
export interface NavbarProps extends React.ComponentPropsWithRef<'nav'> {
  theme?: Theme;
  color?: Color;
}

type NavbarStatic = {
  Brand?: typeof NavbarBrand;
  Nav?: typeof NavbarNav;
  Text?: typeof NavbarText;
};

export type Navbar = React.ForwardRefExoticComponent<
  React.PropsWithRef<NavbarProps>
> &
  NavbarStatic;

export const Navbar: Navbar = React.forwardRef(
  ({ theme, color, className: customClass, children, ...rest }) => {
    const className = classNames(
      'navbar',
      'navbar-expand-lg',
      {
        [`navbar-${theme}`]: theme != null,
        [`bg-${color}`]: color != null
      },
      customClass
    );
    return (
      <nav className={className} {...rest}>
        {children}
      </nav>
    );
  }
);

Navbar.displayName = 'Navbar';

Navbar.Brand = NavbarBrand;
Navbar.Nav = NavbarNav;
Navbar.Text = NavbarText;

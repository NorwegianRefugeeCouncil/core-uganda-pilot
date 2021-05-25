import { classNames } from '@core/ui-toolkit/util/utils';
import * as React from 'react';

type NavbarBrandProps =
  | ({ as?: 'span' } & React.ComponentPropsWithoutRef<'span'>)
  | ({ as: 'a' } & React.ComponentPropsWithoutRef<'a'>);

export const NavbarBrand: React.FC<NavbarBrandProps> = ({
  as = 'span',
  className: customClass,
  children,
  ...rest
}) => {
  const className = classNames('navbar-brand', customClass);
  const Component = as;
  return (
    <Component className={className} {...rest}>
      {children}
    </Component>
  );
};

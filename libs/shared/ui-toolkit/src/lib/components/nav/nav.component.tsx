/* eslint-disable @typescript-eslint/no-empty-interface */
import * as React from 'react';
import { classNames } from '../../helpers/utils';
import { Size, Theme } from '../../helpers/types';

export interface ItemProps extends React.ComponentPropsWithoutRef<'ul'> {}

const Item: React.FC<NavProps> = ({
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


export type NavProps<C extends React.ElementType> & {
  as?: React.ElementType<'ul'>
  size?: Size;
}

const Nav: React.FC<NavProps> = ({
  as = null,
  className: customClass,
  children,
  ...rest
}) => {
  const Component = as ?? 'nav'
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
    <Component className={className} {...rest}>
      <div className="container-fluid">{children}</div>
    </Component>
  );
};



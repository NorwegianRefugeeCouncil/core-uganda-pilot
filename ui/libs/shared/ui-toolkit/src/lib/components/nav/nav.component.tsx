/* eslint-disable @typescript-eslint/no-empty-interface */
import * as React from 'react';
import { classNames } from '@core/ui-toolkit/util/utils';
import { NavItem } from './nav-item.component';
import { NavLink } from './nav-link.component';

export interface NavProps extends React.ComponentPropsWithoutRef<'ul'> {
  variant?: 'tabs' | 'pills';
}

type FillProps =
  | { fill?: true; justified?: never }
  | { fill?: never; justified?: true };

type Nav = React.FC<NavProps & FillProps> & {
  Item: typeof NavItem;
  Link: typeof NavLink;
};

const Nav: Nav = (props) => {
  const {
    variant,
    fill = false,
    justified,
    className: customClass,
    children,
    ...rest
  } = props;
  const className = classNames('nav', customClass, {
    [`nav-${variant}`]: variant != null,
    'nav-fill': fill,
    'nav-justified': justified
  });
  return (
    <ul className={className} {...rest}>
      {children}
    </ul>
  );
};

Nav.displayName = 'Nav';

Nav.Item = NavItem;
Nav.Link = NavLink;

export { Nav };

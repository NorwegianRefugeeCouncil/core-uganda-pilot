import classNames from 'classnames';
import * as React from 'react';
import { uniqueId } from '@ui-helpers/utils';
import { Color, Size } from '@ui-helpers/types';
import Button from '../button/button.component';

// eslint-disable-next-line @typescript-eslint/no-empty-interface
export interface DropdownMenuProps
  extends React.ComponentPropsWithoutRef<'ul'> {
  dark?: true;
}

type AlignProps =
  | { breakEnd?: Size; alignStart?: true; breakStart?: never; alignEnd?: never }
  | {
      breakStart?: Size;
      alignEnd?: true;
      breakEnd?: never;
      alignStart?: never;
    };

type DropdownMenu = React.FC<DropdownMenuProps & AlignProps>;

const DropdownMenu: DropdownMenu = ({
  dark,
  breakEnd,
  breakStart,
  alignStart,
  alignEnd,
  className: customClass,
  children,
  ...rest
}) => {
  const className = classNames('dropdown-menu', customClass, {
    'dropdown-menu-dark': dark,
    [`dropdown-menu-${breakEnd}-end`]: breakEnd,
    [`dropdown-menu-${breakStart}-start`]: breakStart,
    'dropdown-menu-end': alignEnd,
    'dropdown-menu-start': alignStart,
  });
  return (
    <ul className={className} {...rest}>
      {children}
    </ul>
  );
};

export default DropdownMenu;

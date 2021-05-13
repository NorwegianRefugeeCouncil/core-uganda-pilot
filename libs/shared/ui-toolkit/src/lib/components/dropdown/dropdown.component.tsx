import classNames from 'classnames';
import * as React from 'react';
import { uniqueId } from '../../helpers/utils';
import { Color } from '../../helpers/types';
import Button from '../button/button.component';

// eslint-disable-next-line @typescript-eslint/no-empty-interface
export interface DropdownDividerProps
  extends React.ComponentPropsWithoutRef<'hr'> {}

export const DropdownDivider: React.FC<DropdownDividerProps> = (props) => (
  <hr className="dropdown-divider" {...props} />
);

export interface ItemProps extends React.ComponentPropsWithoutRef<'a'> {
  href?: string;
  label: string;
}

export const Item: React.FC<ItemProps> = ({
  href = '#',
  label,
  children,
  ...rest
}) => (
  <a href={href} className="dropdown-item" {...rest}>
    {label ?? children}
  </a>
);

export interface MenuProps extends React.ComponentPropsWithoutRef<'ul'> {
  options?: { type: 'option' | 'divider'; href?: string; label?: string }[];
}

export const Menu: React.FC<MenuProps> = ({
  options = null,
  className: customClass,
  children,
  ...rest
}) => {
  const className = classNames('dropdown-menu', customClass);
  if (options != null) {
    return (
      <ul className={className} {...rest}>
        {options.map((option) => {
          if (option.type === 'option') {
            return (
              <li>
                <Item href={option.href} label={option.label} />
              </li>
            );
          } else if (option.type === 'divider') {
            return (
              <li>
                <DropdownDivider />
              </li>
            );
          } else return null;
        })}
      </ul>
    );
  } else if (children != null) {
    return (
      <ul className={className} {...rest}>
        {children}
      </ul>
    );
  } else return null;
};

export interface DropdownProps extends React.ComponentPropsWithoutRef<'div'> {
  colorTheme?: Color;
  label?: string;
  isOpenInitially?: boolean;
}

const Dropdown: React.FC<DropdownProps> & {
  Menu: typeof Menu;
  Item: typeof Item;
} = ({
  colorTheme = 'primary',
  label = 'Dropdown button',
  isOpenInitially = false,
  className: customClass,
  children,
  ...rest
}) => {
  const [isOpen, setIsOpen] = React.useState(isOpenInitially);
  const className = classNames('dropdown', customClass);
  const id = uniqueId(10);
  return (
    <div className={className} {...rest}>
      <Button
        id={id}
        colorTheme={colorTheme}
        className="dropdown-toggle"
        data-bs-toggle="dropdown"
        aria-expanded={isOpen}
        onPointerDown={() => setIsOpen(!isOpen)}
      >
        {label}
      </Button>
      {isOpen ? children : null}
    </div>
  );
};

Dropdown.displayName = 'Dropdown';
Dropdown.Menu = Menu;
Dropdown.Item = Item;

export default Dropdown;

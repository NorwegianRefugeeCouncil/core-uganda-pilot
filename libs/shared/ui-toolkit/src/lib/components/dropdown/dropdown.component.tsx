import classNames from 'classnames';
import * as React from 'react';
import { UniqueID } from '../../helpers/utils';
import { ThemeColor } from '../../helpers/types';
import { Button } from '../button/button.component';

// eslint-disable-next-line @typescript-eslint/no-empty-interface
export interface DropdownDividerProps
  extends React.ComponentPropsWithoutRef<'hr'> {}

export const DropdownDivider: React.FC<DropdownDividerProps> = (props) => (
  <hr className="dropdown-divider" {...props} />
);

export interface DropdownMenuItemProps
  extends React.ComponentPropsWithoutRef<'a'> {
  href?: string;
  label: string;
}

export const DropdownMenuItem: React.FC<DropdownMenuItemProps> = ({
  href = '#',
  label,
  children,
  ...rest
}) => (
  <a href={href} className="dropdown-item" {...rest}>
    {label ?? children}
  </a>
);

export interface DropdownMenuProps
  extends React.ComponentPropsWithoutRef<'ul'> {
  options?: { type: 'option' | 'divider'; href?: string; label?: string }[];
}

export const DropdownMenu: React.FC<DropdownMenuProps> = ({
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
                <DropdownMenuItem href={option.href} label={option.label} />
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
  theme?: ThemeColor;
  label?: string;
  isOpenInitially?: boolean;
}

export const Dropdown: React.FC<DropdownProps> = ({
  theme = 'primary',
  label = 'Dropdown button',
  isOpenInitially = false,
  className: customClass,
  children,
  ...rest
}) => {
  const [isOpen, setIsOpen] = React.useState(isOpenInitially);
  const className = classNames('dropdown', customClass);
  const id = UniqueID(10);
  return (
    <div className={className} {...rest}>
      <Button
        id={id}
        theme={theme}
        className="dropdown-toggle"
        data-bs-toggle="dropdown"
        aria-expanded={isOpen}
        onPointerDown={() => setIsOpen(!isOpen)}
      >
        {label}
      </Button>
      {children}
    </div>
  );
};

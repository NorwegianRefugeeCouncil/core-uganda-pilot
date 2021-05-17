import * as React from 'react';
import DropdownToggle from './dropdown-toggle.component';
import DropdownMenu from './dropdown-menu.component';
import DropdownItem from './dropdown-item.component';
import DropdownDivider from './dropdown-divider.component';
import Button from '../button/button.component';
import ButtonGroup from '../button-group/button-group.component';
import { classNames } from '@ui-helpers/utils';
import { Color, Direction } from '@ui-helpers/types';
import useDropdown from './use-dropdown';
import { menuPositionFromDir } from './menuPositionFromDir';

export interface DropdownProps extends React.ComponentPropsWithoutRef<'div'> {
  label?: string;
  theme?: Color;
  split?: true;
  dropDir?: Direction;
  onChange?: (value: any) => void;
}

type Dropdown = React.FC<DropdownProps> & {
  Toggle: typeof DropdownToggle;
  Menu: typeof DropdownMenu;
  Item: typeof DropdownItem;
  Divider: typeof DropdownDivider;
};

const Dropdown: Dropdown = ({
  label,
  theme,
  split,
  dropDir = 'down',
  onChange,
  className: customClass,
  children,
  ...rest
}) => {
  const dropDirClass = {
    dropup: dropDir === 'up',
    dropend: dropDir === 'right' || dropDir === 'end',
    dropstart: dropDir === 'left' || dropDir === 'start',
  };

  const {
    menuRef,
    toggleBtnRef,
    menuIsOpen,
    toggleMenu,
    handleChange,
  } = useDropdown(onChange);

  const menuPosition = menuPositionFromDir(dropDir);

  const menu = React.Children.map<React.ReactNode, React.ReactNode>(
    children,
    (child) => {
      if (React.isValidElement(child) && typeof child === typeof DropdownMenu) {
        return React.cloneElement(child, {
          ref: menuRef,
          isVisible: menuIsOpen,
          handleChange,
          position: menuPosition,
        });
      }
    }
  );

  if (split || dropDir != null) {
    const className = classNames(customClass, dropDirClass);
    return (
      <ButtonGroup className={className} {...rest}>
        {split ? (
          <>
            <Button theme={theme} type="button">
              {label}
            </Button>
            <DropdownToggle
              ref={toggleBtnRef}
              split
              theme={theme}
              toggleFn={toggleMenu}
            />
          </>
        ) : (
          <DropdownToggle
            ref={toggleBtnRef}
            theme={theme}
            toggleFn={toggleMenu}
          >
            {label}
          </DropdownToggle>
        )}
        {menu}
      </ButtonGroup>
    );
  } else {
    const className = classNames('dropdown', customClass);
    return (
      <div className={className} {...rest}>
        <DropdownToggle ref={toggleBtnRef} theme={theme} toggleFn={toggleMenu}>
          {label}
        </DropdownToggle>
        {menu}
      </div>
    );
  }
};

Dropdown.displayName = 'Dropdown';

Dropdown.Toggle = DropdownToggle;
Dropdown.Menu = DropdownMenu;
Dropdown.Item = DropdownItem;
Dropdown.Divider = DropdownDivider;

export default Dropdown;

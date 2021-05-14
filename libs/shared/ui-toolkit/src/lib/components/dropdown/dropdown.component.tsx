import * as React from 'react';
import DropdownToggle from './dropdown-toggle.component';
import DropdownMenu from './dropdown-menu.component';
import DropdownItem from './dropdown-item.component';
import DropdownDivider from './dropdown-divider.component';
import Button from '../button/button.component';
import ButtonGroup from '../button-group/button-group.component';
import { classNames } from '@ui-helpers/utils';
import { Color } from '@ui-helpers/types';
import { bsDropdown } from '@ui-helpers/bs-modules';

export interface DropdownProps extends React.ComponentPropsWithoutRef<'div'> {
  label?: string;
  theme?: Color;
  split?: true;
  dropDir?: 'up' | 'right' | 'left' | 'start' | 'end';
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
  dropDir,
  className: customClass,
  children,
  ...rest
}) => {
  const dropDirClass = {
    dropup: dropDir === 'up',
    dropend: dropDir === 'right' || dropDir === 'end',
    dropstart: dropDir === 'left' || dropDir === 'start',
  };

  const [showMenu, setShowMenu] = React.useState(false);
  const toggleMenu = () => setShowMenu(!showMenu);

  const ref = React.useRef();

  React.useEffect(() => {
    const ele = ref.current;
    const dd = new bsDropdown(ref.current, {
      reference: ref.current,
    });
    // if (showMenu) {
    //   dd.show();
    // } else {
    //   dd.hide();
    // }
    // const hideMenu = () => setShowMenu(false);
    // ele.addEventListener('hidden.bs.dropdown', hideMenu);
    // return () => ele.removeEventListener('hidden.bs.dropdown', hideMenu);
  });

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
              ref={ref}
              split
              theme={theme}
              toggleFn={toggleMenu}
            />
          </>
        ) : (
          <DropdownToggle ref={ref} theme={theme} toggleFn={toggleMenu}>
            {label}
          </DropdownToggle>
        )}
        {children}
      </ButtonGroup>
    );
  } else {
    const className = classNames('dropdown', customClass);
    return (
      <div className={className} {...rest}>
        <DropdownToggle ref={ref} theme={theme} toggleFn={toggleMenu}>
          {label}
        </DropdownToggle>
        {children}
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

import * as React from 'react';
import { classNames } from '@ui-helpers/utils';
import Button from '../button/button.component';
import { bsDropdown } from '@ui-helpers/bs-modules';
import { Color } from '@ui-helpers/types';

// eslint-disable-next-line @typescript-eslint/no-empty-interface
export interface DropdownToggleProps
  extends React.ComponentPropsWithRef<'button'> {
  theme?: Color;
  split?: true;
  toggleFn?: () => void;
}

type DropdownToggle = React.FC<DropdownToggleProps>;

const DropdownToggle = ({ theme, split, toggleFn, children, ...rest }, ref) => {
  const className = classNames('dropdown-toggle', {
    'dropdown-toggle-split': split,
  });
  return (
    <Button
      ref={ref}
      type="button"
      theme={theme}
      className={className}
      data-bs-toggle="dropdown"
      aria-expanded="false"
      onPointerDown={toggleFn}
      {...rest}
    >
      {children}
    </Button>
  );
};

export default React.forwardRef<HTMLButtonElement, DropdownToggleProps>(
  DropdownToggle
);

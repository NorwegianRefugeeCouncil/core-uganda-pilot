import * as React from 'react';
import { classNames } from '@core/shared/ui-toolkit/util/utils';
import { Button } from '../button/button.component';
import { Color } from '@core/shared/ui-toolkit/util/types';

// eslint-disable-next-line @typescript-eslint/no-empty-interface
export interface DropdownToggleProps
  extends React.ComponentPropsWithRef<'button'> {
  theme?: Color;
  split?: true;
  toggleFn?: (event: React.PointerEvent) => void;
}

type DropdownToggle = React.FC<DropdownToggleProps>;

export const DropdownToggle = React.forwardRef<
  HTMLButtonElement,
  DropdownToggleProps
>(({ theme, split, toggleFn, children, ...rest }, ref) => {
  const className = classNames('dropdown-toggle', {
    'dropdown-toggle-split': split,
  });
  return (
    <Button
      ref={ref}
      type="button"
      theme={theme}
      className={className}
      onPointerDown={toggleFn}
      {...rest}
    >
      {children}
    </Button>
  );
});

DropdownToggle.displayName = 'DropdownToggle';

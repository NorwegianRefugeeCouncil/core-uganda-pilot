import classNames from 'classnames';
import * as React from 'react';
import { Size } from '@core/shared/ui-toolkit/util/types';

// eslint-disable-next-line @typescript-eslint/no-empty-interface
export interface BasicDropdownMenuProps
  extends React.ComponentPropsWithRef<'ul'> {
  isVisible?: boolean;
  handleChange: (value: any) => void;
  position: React.CSSProperties;
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

type DropdownMenuProps = BasicDropdownMenuProps & AlignProps;

const DropdownMenu = React.forwardRef<HTMLUListElement, DropdownMenuProps>(
  (
    {
      isVisible = false,
      dark,
      breakEnd,
      breakStart,
      alignStart,
      alignEnd,
      position,
      handleChange,
      className: customClass,
      children,
      ...rest
    },
    ref
  ) => {
    const className = classNames('dropdown-menu', customClass, {
      'dropdown-menu-dark': dark,
      [`dropdown-menu-${breakEnd}-end`]: breakEnd,
      [`dropdown-menu-${breakStart}-start`]: breakStart,
      'dropdown-menu-end': alignEnd,
      'dropdown-menu-start': alignStart,
      show: isVisible,
    });
    return (
      <ul ref={ref} className={className} style={position} {...rest}>
        {React.Children.map<React.ReactNode, React.ReactNode>(
          children,
          (child) => {
            if (React.isValidElement(child)) {
              return React.cloneElement(child, { handleChange });
            }
          }
        )}
      </ul>
    );
  }
);

DropdownMenu.displayName = 'DropdownMenu';

export { DropdownMenu };

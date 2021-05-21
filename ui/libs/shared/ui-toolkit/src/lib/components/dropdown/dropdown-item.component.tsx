import * as React from 'react';
import { classNames } from '@core/ui-toolkit/util/utils';

// eslint-disable-next-line @typescript-eslint/no-empty-interface
export interface DropdownItemProps
  extends React.ComponentPropsWithoutRef<'li'> {
  href?: string;
  isTextOnly?: true;
  value?: any;
  handleChange?: (value: any) => void;
}

type ActiveProps =
  | { active?: true; disabled?: never }
  | { active?: never; disabled?: true };

type DropdownItem = React.FC<DropdownItemProps & ActiveProps>;

export const DropdownItem: DropdownItem = ({
  href,
  isTextOnly,
  value,
  active,
  disabled,
  handleChange,
  className: customClass,
  children,
  ...rest
}) => {
  const className = classNames('dropdown-item', customClass, {
    active,
    disabled,
    'dropdown-item-text': isTextOnly,
  });
  return (
    <li {...rest}>
      <a
        className={className}
        href={href}
        onPointerDown={(value) => handleChange(value)}
      >
        {children}
      </a>
    </li>
  );
};

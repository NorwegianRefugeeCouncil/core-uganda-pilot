import * as React from 'react';
import { classNames, uniqueId } from '@ui-helpers/utils';

// eslint-disable-next-line @typescript-eslint/no-empty-interface
export interface DropdownItemProps
  extends React.ComponentPropsWithoutRef<'li'> {
  href?: string;
  text?: true;
}

type ActiveProps =
  | { active?: true; disabled?: never }
  | { active?: never; disabled?: true };

type DropdownItem = React.FC<DropdownItemProps & ActiveProps>;

const DropdownItem: DropdownItem = ({
  href,
  text,
  active,
  disabled,
  className: customClass,
  children,
  ...rest
}) => {
  const className = classNames('dropdown-item', customClass, {
    active,
    disabled,
    'dropdown-item-text': text,
  });
  return (
    <li {...rest}>
      <a className={className} href={href}>
        {children}
      </a>
    </li>
  );
};

export default DropdownItem;

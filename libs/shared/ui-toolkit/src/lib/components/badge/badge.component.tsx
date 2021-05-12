import * as React from 'react';
import classNames from 'classnames';
import { ThemeColor } from '../../helpers/types';

export interface BadgeProps extends React.ComponentPropsWithoutRef<'span'> {
  theme?: ThemeColor;
  pill?: boolean;
}

export const Badge: React.FC<BadgeProps> = ({
  className: customClass,
  children,
  theme = 'primary',
  pill = false,
  ...rest
}) => {
  const className = classNames(
    'badge',
    `bg-${theme}`,
    {
      'text-dark': theme === 'light' || theme === 'warning',
      'text-light': theme === 'dark',
      'rounded-pill': pill,
    },
    customClass
  );
  return (
    <span className={className} {...rest}>
      {children}
    </span>
  );
};

export default Badge;

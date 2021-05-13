import * as React from 'react';
import classNames from 'classnames';
import { Color, Size } from '../../helpers/types';
import { color } from '@storybook/addon-knobs';

export interface ButtonProps extends React.ComponentPropsWithRef<'button'> {
  colorTheme?: Color | 'link';
  size?: 'sm' | 'lg';
  outline?: boolean;
}

type Button = React.FC<ButtonProps>;

const Button: Button = (props: ButtonProps) => {
  const {
    colorTheme = 'primary',
    size,
    outline = false,
    className: customClass,
    children,
    ...rest
  } = props;
  const className = classNames('btn', customClass, {
    [`btn-${colorTheme}`]: colorTheme && !outline,
    [`btn-outline-${colorTheme}`]: outline,
    [`btn-${size}`]: size != null,
  });
  return (
    <button {...rest} className={className}>
      {children}
    </button>
  );
};

Button.displayName = 'Button';

export default Button;

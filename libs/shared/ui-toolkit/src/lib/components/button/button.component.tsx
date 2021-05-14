import * as React from 'react';
import classNames from 'classnames';
import { Color, Size } from '../../helpers/types';

export interface ButtonProps extends React.ComponentPropsWithRef<'button'> {
  theme?: Color | 'link';
  size?: 'sm' | 'lg';
  outline?: boolean;
}

const Button = (props, ref) => {
  const {
    theme = 'primary',
    size,
    outline = false,
    className: customClass,
    children,
    ...rest
  } = props;
  const className = classNames('btn', customClass, {
    [`btn-${theme}`]: theme && !outline,
    [`btn-outline-${theme}`]: outline,
    [`btn-${size}`]: size != null,
  });
  return (
    <button ref={ref} className={className} {...rest}>
      {children}
    </button>
  );
};

Button.displayName = 'Button';

export default React.forwardRef<HTMLButtonElement, ButtonProps>(Button);

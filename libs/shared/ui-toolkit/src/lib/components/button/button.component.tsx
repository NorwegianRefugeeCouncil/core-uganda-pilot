import React, { ButtonHTMLAttributes, FunctionComponent } from 'react';
import classNames from 'classnames';

type ButtonProps = ButtonHTMLAttributes<HTMLButtonElement> & {
  theme?:
    | 'primary'
    | 'secondary'
    | 'danger'
    | 'success'
    | 'warning'
    | 'info'
    | 'light'
    | 'dark'
    | 'link';
  size?: 'sm' | 'lg';
  outline?: boolean;
};

export const Button: FunctionComponent<ButtonProps> = (props: ButtonProps) => {
  const {
    theme = 'primary',
    size,
    outline,
    className,
    children,
    ...otherProps
  } = props;
  const classes: string[] = [];
  classes.push('btn');
  if (theme) {
    classes.push('btn-' + (outline ? 'outline-' : '') + theme);
  }
  if (size) {
    classes.push('btn-' + size);
  }
  const btnClass = classNames(className, classes);
  return (
    <button {...otherProps} className={btnClass}>
      {children}
    </button>
  );
};

type CloseButtonProps = ButtonHTMLAttributes<HTMLButtonElement> & {
  size?: 'sm' | 'lg';
};

export const CloseButton = (() => {
  const cmp: FunctionComponent<CloseButtonProps> = ({
    size,
    className,
    children,
    ...props
  }: ButtonProps) => {
    className = classNames(className, {
      btn: true,
      'btn-close': true,
      'btn-sm': size === 'sm',
      'btn-lg': size === 'lg',
    });
    return (
      <button className={className} {...props}>
        {children}
      </button>
    );
  };
  cmp.defaultProps = {
    type: 'button',
    'aria-label': 'Close',
  };
  cmp.displayName = 'CloseButton';
  return cmp;
})();

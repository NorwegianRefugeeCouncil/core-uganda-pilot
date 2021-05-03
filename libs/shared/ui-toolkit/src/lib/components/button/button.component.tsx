import React, {
  ButtonHTMLAttributes,
} from 'react';
import classNames from 'classnames';

type ButtonProps = ButtonHTMLAttributes<HTMLButtonElement> & {
  kind?: 'primary' | 'secondary' | 'danger' | 'success' | 'warning' | 'info' | 'light' | 'dark' | 'link'
  size?: 'sm' | 'lg'
  outline?: boolean
}

export const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  (props: ButtonProps, ref) => {
    const {kind, size, outline, className, children, ...otherProps} = props;
    const classes: string[] = [];
    classes.push('btn');
    if (kind) {
      classes.push('btn-' + (outline ? 'outline-' : '') + kind);
    }
    if (size) {
      classes.push('btn-' + size);
    }
    return (<button ref={ref} {...otherProps} className={classNames(className, ...classes)}>{children}</button>);
  });

type CloseButtonProps = ButtonHTMLAttributes<HTMLButtonElement> & {
  size?: 'sm' | 'lg'
}

export const CloseButton = (() => {
  const cmp = React.forwardRef<HTMLButtonElement, CloseButtonProps>(
    ({size, className, children, ...props}: ButtonProps, ref) => {
      className = classNames(className, 'btn-close');
      if (size) {
        className = classNames(className, 'btn-' + size);
      }
      return <button
        ref={ref}
        className={className}
        {...props}>{children}
      </button>;
    });
  cmp.defaultProps = {
    type: 'button',
    'aria-label': 'Close'
  };
  cmp.displayName = 'CloseButton';
  return cmp;
})();

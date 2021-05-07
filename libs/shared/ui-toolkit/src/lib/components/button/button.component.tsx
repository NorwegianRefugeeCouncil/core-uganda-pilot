import * as React from 'react';
import classNames from 'classnames';

interface ButtonProps extends React.ComponentPropsWithRef<'button'> {
  kind?:
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
  type?: 'submit' | 'button';
  outline?: boolean;
}

export const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  (
    {
      type = 'submit',
      kind = 'primary',
      outline = false,
      size,
      className,
      children,
      ...baseProps
    },
    ref
  ) => {
    const classes: string[] = [];
    classes.push('btn');
    if (kind) {
      classes.push('btn-' + (outline ? 'outline-' : '') + kind);
    }
    if (size) {
      classes.push('btn-' + size);
    }
    const btnClass = classNames(className, classes);
    return (
      <button ref={ref} type={type} {...baseProps} className={btnClass}>
        {children}
      </button>
    );
  }
);

export const CloseButton = React.forwardRef<
  HTMLButtonElement,
  ButtonProps & { white?: boolean }
>(({ size, white = false, className, ...baseProps }, ref) => {
  const classes = classNames(className, {
    'btn-close': true,
    'btn-close-white': white,
    'btn-sm': size === 'sm',
    'btn-lg': size === 'lg',
  });
  return (
    <button
      ref={ref}
      {...baseProps}
      type="button"
      aria-label="Close"
      className={classes}
    />
  );
});

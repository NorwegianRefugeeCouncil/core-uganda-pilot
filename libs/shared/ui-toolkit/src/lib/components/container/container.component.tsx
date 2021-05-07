import * as React from 'react';
import classNames from 'classnames';

export interface ContainerProps extends React.ComponentProps<'div'> {
  /**
   * If `true`, container will center its children
   * regardless of their width.
   */
  size?: 'sm' | 'md' | 'lg' | 'xl' | 'xxl' | 'fluid';
  centerContent?: boolean;
}

/**
 * Layout component used to wrap app or website content
 *
 * It sets `margin-left` and `margin-right` to `auto`,
 * to keep its content centered.
 *
 * It also sets a default max-width of `60ch` (60 characters).
 */
export const Container = React.forwardRef<HTMLDivElement, ContainerProps>(
  (
    {
      size,
      centerContent = false,
      className,
      children,
      ...otherProps
    }: ContainerProps,
    ref
  ) => {
    const containerClass = size ? `container-${size}` : 'container';
    return (
      <div
        ref={ref}
        className={classNames(className, containerClass)}
        {...otherProps}
        style={
          centerContent && {
            display: 'flex',
            justifyContent: 'center',
            alignItems: 'center',
          }
        }
      >
        {children}
      </div>
    );
  }
);

import * as React from 'react';
import classNames from 'classnames';

export interface ContainerProps extends React.ComponentPropsWithRef<'div'> {
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
const Container = React.forwardRef<HTMLDivElement, ContainerProps>(
  (
    { size, centerContent = false, className: customClass, children, ...rest },
    ref
  ) => {
    const className = classNames('container', customClass, {
      [`container-${size}`]: size != null,
    });
    return (
      <div
        ref={ref}
        className={className}
        {...rest}
        style={
          centerContent
            ? {
                display: 'flex',
                justifyContent: 'center',
                alignItems: 'center',
              }
            : {}
        }
      >
        {children}
      </div>
    );
  }
);

export default Container;

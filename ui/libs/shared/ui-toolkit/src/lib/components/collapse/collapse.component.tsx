import * as React from 'react';
import classNames from 'classnames';

export type CollapseProps = React.ComponentPropsWithRef<'div'> & {
  show?: true;
};

const Collapse = React.forwardRef<HTMLDivElement, CollapseProps>(
  (props, ref) => {
    const { show, className: customClass, children, ...rest } = props;
    const classeName = classNames(customClass, 'collapse', { show });
    return (
      <div {...rest} className={classeName}>
        {children}
      </div>
    );
  }
);

Collapse.displayName = 'Collapse';

export { Collapse };

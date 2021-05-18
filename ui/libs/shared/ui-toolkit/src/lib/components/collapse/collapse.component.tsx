import * as React from 'react';
import classNames from 'classnames';
import { TransitionState } from '@ui-helpers/transition';

export type CollapseProps = React.ComponentPropsWithRef<'div'> & {
  show?: boolean;
};

export const Collapse = React.forwardRef<HTMLDivElement, CollapseProps>(
  (props, ref) => {
    const { show = true, className: customClass, children, ...rest } = props;
    const classeName = classNames(customClass, 'collapse', { show });
    return (
      <div {...rest} className={classeName}>
        {children}
      </div>
    );
  }
);

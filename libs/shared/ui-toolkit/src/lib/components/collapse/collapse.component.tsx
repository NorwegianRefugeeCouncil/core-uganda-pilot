import * as React from 'react';
import classNames from 'classnames';

export type CollapseProps = React.ComponentPropsWithRef<'div'> & {
  show?: boolean;
};

export const Collapse: React.FC<CollapseProps> = (props) => {
  const classes = classNames(props.className, 'collapse', {
    show: props.show,
  });
  return (
    <div {...props} className={classes}>
      {props.children}
    </div>
  );
};

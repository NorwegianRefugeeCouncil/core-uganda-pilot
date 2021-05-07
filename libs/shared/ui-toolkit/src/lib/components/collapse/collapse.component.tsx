import { FunctionComponent, HTMLAttributes } from 'react';
import classNames from 'classnames';

export type CollapseProps = HTMLAttributes<HTMLDivElement> & {
  show?: boolean;
};

export const Collapse: FunctionComponent<CollapseProps> = (props) => {
  const classes = classNames(props.className, 'collapse', {
    show: props.show,
  });
  return (
    <div {...props} className={classes}>
      {props.children}
    </div>
  );
};

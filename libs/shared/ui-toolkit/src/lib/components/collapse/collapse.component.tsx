import {
  forwardRef,
  FunctionComponent,
  HTMLAttributes,
  LabelHTMLAttributes,
} from 'react';
import classNames from 'classnames';

export type CollapseProps = HTMLAttributes<HTMLDivElement> & {
  show?: boolean;
};

export const Collapse = forwardRef<HTMLDivElement, CollapseProps>(
  (props, ref) => {
    const classes = classNames(props.className, 'collapse', {
      show: props.show,
    });
    return (
      <div {...props} ref={ref} className={classes}>
        {props.children}
      </div>
    );
  }
);

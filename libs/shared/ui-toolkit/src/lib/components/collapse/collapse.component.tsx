import { forwardRef, FunctionComponent, HTMLAttributes, LabelHTMLAttributes } from 'react';
import { addClasses } from '../../utils/utils';

type CollapseProps = HTMLAttributes<HTMLDivElement> & {
  show?: boolean
}

export const Collapse = forwardRef<HTMLDivElement, CollapseProps>(
  (props, ref) => {
    const classes: string[] = [];
    classes.push('collapse');
    if (props.show) {
      classes.push('show');
    }
    return (<div {...props} ref={ref} className={addClasses(props.className, ...classes)}>{props.children}</div>);
  });

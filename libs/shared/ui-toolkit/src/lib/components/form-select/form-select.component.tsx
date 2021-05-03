import {
  forwardRef,
  FunctionComponent, InputHTMLAttributes, LabelHTMLAttributes, SelectHTMLAttributes
} from 'react';
import { addClasses } from '../../utils/utils';
import { CollapseProps } from '../collapse/collapse.component';

type SelectProps = SelectHTMLAttributes<HTMLSelectElement>

export const FormSelect = forwardRef<HTMLSelectElement, SelectProps>(
  (props, ref) => {
    return (<select {...props} ref={ref} className={addClasses(props.className, 'form-select')} />);
  });

import { Color } from '@ui-helpers/types';
import { classNames } from '@ui-helpers/utils';
import * as React from 'react';

export interface ProgessBarProps extends React.ComponentPropsWithRef<'div'> {
  label?: string;
  striped?: boolean;
  animated?: boolean;
  color?: Color;
  progress?: number;
}

const ProgressBar: React.FC<ProgessBarProps> = ({
  label,
  striped = false,
  animated = false,
  color,
  progress = 0,
  className: customClass,
  children,
  ...rest
}) => {
  const className = classNames(
    'progress-bar',
    {
      'progress-bar-striped': striped,
      'progress-bar-animated': animated,
    },
    customClass
  );
  return (
    <div
      className={className}
      role="progressbar"
      style={{ width: progress + '%' }}
      {...rest}
      aria-valuenow={progress}
    >
      {children}
    </div>
  );
};

export default ProgressBar;

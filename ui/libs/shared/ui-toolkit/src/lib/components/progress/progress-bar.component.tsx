import { Color } from '@ui-helpers/types';
import { classNames } from '@ui-helpers/utils';
import * as React from 'react';

export interface ProgessBarProps extends React.ComponentPropsWithRef<'div'> {
  label?: string;
  striped?: boolean;
  animated?: boolean;
  theme?: Color;
  progress?: number;
}

export const ProgressBar: React.FC<ProgessBarProps> = ({
  label,
  striped = false,
  animated = false,
  theme,
  progress = 0,
  className: customClass,
  children,
  ...rest
}) => {
  const className = classNames(
    'progress-bar',
    {
      'progress-bar-striped': striped && !animated,
      'progress-bar-striped progress-bar-animated': animated,
      [`bg-${theme}`]: theme != null,
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
      aria-valuemin={0}
      aria-valuemax={100}
    >
      {children}
    </div>
  );
};

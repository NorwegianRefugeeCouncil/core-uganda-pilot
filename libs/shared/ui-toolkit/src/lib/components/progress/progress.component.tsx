import { Color } from '@ui-helpers/types';
import { classNames } from '@ui-helpers/utils';
import * as React from 'react';
import ProgressBar from './progress-bar.component';

export interface ProgressProps extends React.ComponentPropsWithRef<'div'> {
  showValue?: boolean;
  progress?: number;
  theme?: Color;
  striped?: boolean;
  animated?: boolean;
  height?: number;
  children?: React.ReactText;
}

type Progress = React.FC<ProgressProps> & {
  Bar: typeof ProgressBar;
};

const Progress: Progress = ({
  showValue = false,
  progress = 0,
  theme,
  striped = false,
  animated = false,
  height = 20,
  children: label,
  className: customClass,
  children,
  ...rest
}) => {
  if (progress < 0 || progress > 100)
    throw new RangeError('"progress" prop should be in range 0 to 100');
  const className = classNames('progress', customClass);
  return (
    <div
      className={className}
      style={{ height: height ? height : 'initial' }}
      {...rest}
    >
      {children ?? (
        <ProgressBar
          progress={progress}
          theme={theme}
          striped={striped}
          animated={animated}
        >
          {showValue ? `${progress}%` : null}
        </ProgressBar>
      )}
    </div>
  );
};

Progress.displayName = 'Progress';

Progress.Bar = ProgressBar;

export default Progress;

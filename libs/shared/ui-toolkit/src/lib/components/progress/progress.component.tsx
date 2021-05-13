import { Color } from '@ui-helpers/types';
import { classNames } from '@ui-helpers/utils';
import * as React from 'react';
import ProgressBar from './progress-bar.component';

export interface ProgessProps extends React.ComponentPropsWithRef<'div'> {
  labels?: string[];
  showValue?: boolean;
  progress?: number;
  color?: Color;
  height?: number;
  children?: React.ReactText;
}

const Progress: React.FC<ProgessProps> = ({
  labels = null,
  showValue = false,
  progress = 0,
  color,
  height = 20,
  children: label,
  className: customClass,
  children,
  ...rest
}) => {
  const className = classNames('progress');
  return (
    <>
      <div
        className={className}
        style={{ height: height ? height : 'initial' }}
        {...rest}
      >
        <ProgressBar progress={progress} color={color}>
          {children}
        </ProgressBar>
      </div>
      <div className="d-flex justify-content-between">
        {labels ? labels.map((l) => <span>{l}</span>) : null}
      </div>
    </>
  );
};

export default Progress;

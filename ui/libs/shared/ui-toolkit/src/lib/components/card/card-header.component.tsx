import * as React from 'react';
import { classNames } from '@ui-helpers/utils';

// eslint-disable-next-line @typescript-eslint/no-empty-interface
export interface CardHeaderProps extends React.ComponentPropsWithRef<'div'> {}

export const CardHeader: React.FC<CardHeaderProps> = (props) => {
  return (
    <div {...props} className={classNames(props.className, 'card-header')}>
      {props.children}
    </div>
  );
};

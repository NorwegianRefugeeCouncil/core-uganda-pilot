import * as React from 'react';
import { classNames } from '@core/ui-toolkit/util/utils';

// eslint-disable-next-line @typescript-eslint/no-empty-interface
export interface CardTextProps extends React.ComponentPropsWithRef<'p'> {}

export const CardText: React.FC<CardTextProps> = (props) => {
  return (
    <p {...props} className={classNames(props.className, 'card-text')}>
      {props.children}
    </p>
  );
};

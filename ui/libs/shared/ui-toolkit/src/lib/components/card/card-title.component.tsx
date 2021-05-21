import * as React from 'react';
import { classNames } from '@core/ui-toolkit/util/utils';

// eslint-disable-next-line @typescript-eslint/no-empty-interface
export interface CardTitleProps extends React.ComponentPropsWithRef<'h5'> {}

export const CardTitle: React.FC<CardTitleProps> = (props) => {
  return (
    <h5 {...props} className={classNames(props.className, 'card-title')}>
      {props.children}
    </h5>
  );
};

import * as React from 'react';
import { classNames } from '@core/shared/ui-toolkit/util/utils';

// eslint-disable-next-line @typescript-eslint/no-empty-interface
export interface CardFooterProps extends React.ComponentPropsWithRef<'div'> {}

export const CardFooter: React.FC<CardFooterProps> = (props) => {
  return (
    <div {...props} className={classNames(props.className, 'card-footer')}>
      {props.children}
    </div>
  );
};

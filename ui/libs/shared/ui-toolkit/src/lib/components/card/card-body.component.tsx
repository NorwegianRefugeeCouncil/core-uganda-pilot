import * as React from 'react';
import { classNames } from '@ui-helpers/utils';

// eslint-disable-next-line @typescript-eslint/no-empty-interface
export interface CardBodyProps extends React.ComponentPropsWithRef<'div'> {}

export const CardBody: React.FC<CardBodyProps> = (props) => {
  return (
    <div {...props} className={classNames(props.className, 'card-body')}>
      {props.children}
    </div>
  );
};

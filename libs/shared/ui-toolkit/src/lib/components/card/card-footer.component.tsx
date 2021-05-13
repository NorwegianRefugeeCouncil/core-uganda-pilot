import * as React from 'react';
import { classNames } from '@ui-helpers/utils';

// eslint-disable-next-line @typescript-eslint/no-empty-interface
export interface CardFooterProps extends React.ComponentPropsWithRef<'div'> {}

const CardFooter: React.FC<CardFooterProps> = (props) => {
  return (
    <div {...props} className={classNames(props.className, 'card-footer')}>
      {props.children}
    </div>
  );
};

export default CardFooter;

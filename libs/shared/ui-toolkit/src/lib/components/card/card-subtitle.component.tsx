import * as React from 'react';
import { classNames } from '@ui-helpers/utils';

// eslint-disable-next-line @typescript-eslint/no-empty-interface
export interface CartSubtitleProps extends React.ComponentPropsWithRef<'h6'> {}

const CartSubtitle: React.FC<CartSubtitleProps> = (props) => {
  return (
    <h6 {...props} className={classNames(props.className, 'card-subtitle')}>
      {props.children}
    </h6>
  );
};

export default CartSubtitle;

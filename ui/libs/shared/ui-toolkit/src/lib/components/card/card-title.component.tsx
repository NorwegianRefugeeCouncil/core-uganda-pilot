import * as React from 'react';
import { classNames } from '@ui-helpers/utils';

// eslint-disable-next-line @typescript-eslint/no-empty-interface
export interface CartTitleProps extends React.ComponentPropsWithRef<'h5'> {}

const CartTitle: React.FC<CartTitleProps> = (props) => {
  return (
    <h5 {...props} className={classNames(props.className, 'card-title')}>
      {props.children}
    </h5>
  );
};

export default CartTitle;

import * as React from 'react';
import { classNames } from '@ui-helpers/utils';

// eslint-disable-next-line @typescript-eslint/no-empty-interface
export interface CardTextProps extends React.ComponentPropsWithRef<'p'> {}

const CardText: React.FC<CardTextProps> = (props) => {
  return (
    <p {...props} className={classNames(props.className, 'card-text')}>
      {props.children}
    </p>
  );
};

export default CardText;

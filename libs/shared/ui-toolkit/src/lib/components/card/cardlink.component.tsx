import * as React from 'react';
import { classNames } from '@ui-helpers/utils';

// eslint-disable-next-line @typescript-eslint/no-empty-interface
export interface CardLinkProps extends React.ComponentPropsWithRef<'a'> {}

const CardLink: React.FC<CardLinkProps> = (props) => {
  return (
    <a {...props} className={classNames(props.className, 'card-link')}>
      {props.children}
    </a>
  );
};

export default CardLink;

import * as React from 'react';
import { classNames } from '@core/ui-toolkit/util/utils';

// eslint-disable-next-line @typescript-eslint/no-empty-interface
export interface CardLinkProps extends React.ComponentPropsWithRef<'a'> {}

export const CardLink: React.FC<CardLinkProps> = (props) => {
  return (
    <a {...props} className={classNames(props.className, 'card-link')}>
      {props.children}
    </a>
  );
};

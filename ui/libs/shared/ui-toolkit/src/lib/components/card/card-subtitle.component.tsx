import * as React from 'react';
import { classNames } from '@core/ui-toolkit/util/utils';

// eslint-disable-next-line @typescript-eslint/no-empty-interface
export interface CardSubtitleProps extends React.ComponentPropsWithRef<'h6'> {}

export const CardSubtitle: React.FC<CardSubtitleProps> = (props) => {
  return (
    <h6 {...props} className={classNames(props.className, 'card-subtitle')}>
      {props.children}
    </h6>
  );
};

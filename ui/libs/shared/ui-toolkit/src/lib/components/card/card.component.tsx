/* eslint-disable @typescript-eslint/no-empty-interface */
import * as React from 'react';
import { classNames } from '@core/ui-toolkit/util/utils';
import { PropsWithoutRef, RefAttributes } from 'react';

export interface CardProps extends React.ComponentPropsWithoutRef<'div'> {
}

type Card = React.ForwardRefExoticComponent<PropsWithoutRef<CardProps> & RefAttributes<HTMLDivElement>>

const Card: Card = React.forwardRef<HTMLDivElement, CardProps>((
  {
    className: customClass,
    children,
    ...rest
  }, ref) => {
  const className = classNames('card', customClass);
  return (
    <div ref={ref} {...rest} className={className}>
      {children}
    </div>
  );
});

Card.displayName = 'Card';

export { Card };

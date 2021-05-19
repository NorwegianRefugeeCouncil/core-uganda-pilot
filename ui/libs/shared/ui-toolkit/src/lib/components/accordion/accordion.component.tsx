import * as React from 'react';
import classNames from 'classnames';
import { AccordionItem } from './accordion-item.component';

export interface AccordionProps extends React.ComponentPropsWithoutRef<'div'> {
  isFlush?: boolean;
  stayOpen?: boolean;
}

type Accordion = React.FC<AccordionProps> & {
  Item: typeof AccordionItem;
};

export const Accordion: Accordion = ({
  isFlush = false,
  stayOpen = false,
  className: customClass,
  children,
  ...rest
}) => {
  const accordionClasses = classNames('accordion', customClass, {
    'accordion-flush': isFlush,
  });
  return (
    <div className={accordionClasses} {...rest}>
      {children}
    </div>
  );
};

Accordion.displayName = 'Accordion';

Accordion.Item = AccordionItem;

import * as React from 'react';
import classNames from 'classnames';
import { useAccordion } from './use-accordion';
import { AccordionItem } from './accordion-item.component';
import { AccordionContext } from './accordion-context';

export interface AccordionProps extends React.ComponentPropsWithRef<'div'> {
  activeId: string;
  defaultActiveId?: string;
  flush?: true;
  onSelection?: (key: string) => void;
}

export type AccordionStatic = {
  Item?: typeof AccordionItem;
};

export type Accordion = React.ForwardRefExoticComponent<
  React.PropsWithRef<AccordionProps>
> &
  AccordionStatic;

export const Accordion: Accordion = React.forwardRef<
  HTMLDivElement,
  AccordionProps
>(
  (
    {
      activeId: id,
      defaultActiveId,
      flush,
      onSelection: onPointerDown = null,
      className: customClass,
      children,
      ...rest
    },
    ref
  ) => {
    const { activeKey, handlePointerDown } = useAccordion({
      id: defaultActiveId ?? id,
      onPointerDown,
    });

    const className = classNames(
      'accordion',
      {
        'accordion-flush': flush,
      },
      customClass
    );
    return (
      <AccordionContext.Provider value={{ activeKey, handlePointerDown }}>
        <div ref={ref} className={className} {...rest}>
          {children}
        </div>
      </AccordionContext.Provider>
    );
  }
);

Accordion.displayName = 'Accordion';

Accordion.Item = AccordionItem;

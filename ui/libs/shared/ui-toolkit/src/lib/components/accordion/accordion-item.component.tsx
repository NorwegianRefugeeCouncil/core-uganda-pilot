import * as React from 'react';
import classNames from 'classnames';
import { AccordionHeader } from './accordion-header.component';
import { AccordionContext } from './accordion-context';
import { AccordionBody } from './accordion-body.component';
import { AccordionCollapse } from './accordion-collapse.component';

export interface AccordionItemProps extends React.ComponentPropsWithRef<'div'> {
  id: string;
  header: string;
  body: string | JSX.Element;
  open?: true;
}

export const AccordionItem = React.forwardRef<
  HTMLDivElement,
  AccordionItemProps
>(({ id, header, body, open, className: customClass, ...rest }, ref) => {
  const { activeKey } = React.useContext(AccordionContext);
  const itemClassName = classNames('accordion-item', customClass);
  const collapseClassName = classNames('accordion-collapse', 'collapse', {
    show: open || id === activeKey,
  });
  return (
    <div ref={ref} className={itemClassName} {...rest}>
      <AccordionHeader id={id}>{header}</AccordionHeader>
      <AccordionCollapse id={id} className={collapseClassName}>
        <AccordionBody>{body}</AccordionBody>
      </AccordionCollapse>
    </div>
  );
});

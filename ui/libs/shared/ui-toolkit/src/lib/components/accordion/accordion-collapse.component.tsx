import * as React from 'react';
import { AccordionContext } from './accordion-context';
import { classNames } from '@core/ui-toolkit/util/utils';

export interface AccordionCollapseProps
  extends React.ComponentPropsWithRef<'div'> {
  id: string;
}

export const AccordionCollapse = React.forwardRef<
  HTMLHeadingElement,
  AccordionCollapseProps
>(({ id, ...rest }, ref) => {
  const { activeKey, handlePointerDown } = React.useContext(AccordionContext);

  const className = classNames('accordion-button', {
    collapsed: id !== activeKey,
  });

  return <div ref={ref} id={id} className="accordion-header" {...rest} />;
});

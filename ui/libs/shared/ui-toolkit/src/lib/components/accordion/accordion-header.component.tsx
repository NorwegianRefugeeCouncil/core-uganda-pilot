import * as React from 'react';
import { AccordionContext } from './accordion-context';
import { classNames } from '@core/shared/ui-toolkit/util/utils';

export interface AccordionHeaderProps
  extends React.ComponentPropsWithRef<'h2'> {
  id: string;
}

export const AccordionHeader = React.forwardRef<
  HTMLHeadingElement,
  AccordionHeaderProps
>(({ id, children, ...rest }, ref) => {
  const { activeKey, handlePointerDown } = React.useContext(AccordionContext);

  const className = classNames('accordion-button', {
    collapsed: id !== activeKey,
  });

  return (
    <h2 ref={ref} id={id} className="accordion-header" {...rest}>
      <button
        className={className}
        type="button"
        onPointerDown={() => handlePointerDown(id)}
      >
        {children}
      </button>
    </h2>
  );
});

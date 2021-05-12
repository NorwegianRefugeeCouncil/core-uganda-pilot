import * as React from 'react';
import classNames from 'classnames';
import { uniqueId } from '../../helpers/utils';

interface ItemProps extends React.ComponentPropsWithoutRef<'div'> {
  title: string;
  body: string | HTMLElement;
  isCollapsed?: boolean;
}

const Item: React.FC<ItemProps> = ({
  title,
  body,
  isCollapsed = true,
  ...baseProps
}) => {
  const id = uniqueId(6);
  const headerId = 'header-' + id;
  const collapseId = 'collapse-' + id;
  const buttonClass = classNames('accordion-button', {
    collapsed: isCollapsed,
  });
  const collapseClass = classNames('accordion-collapse collapse', {
    show: !isCollapsed,
  });
  return (
    <>
      <h2 className="accordion-header" id={headerId}>
        <button
          className={buttonClass}
          type="button"
          aria-expanded={isCollapsed}
          aria-controls={collapseId}
        >
          {title}
        </button>
      </h2>
      <div id={collapseId} className={collapseClass} aria-labelledby={headerId}>
        <div className="accordion-body">{body}</div>
      </div>
    </>
  );
};

interface AccordionProps extends React.ComponentPropsWithoutRef<'div'> {
  isFlush?: boolean;
  stayOpen?: boolean;
}

const Accordion: React.FC<AccordionProps> = ({
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

Accordion.Item = Item;

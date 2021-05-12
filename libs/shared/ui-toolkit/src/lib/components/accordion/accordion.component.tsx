import * as React from 'react';
import classNames from 'classnames';
import { ProgressPlugin } from 'webpack';

interface AccordionProps extends React.ComponentPropsWithoutRef<'div'> {
  isFlush?: boolean;
  stayOpen?: boolean;
}

const Accordion: React.FC<AccordionProps> = ({
  isFlush = false,
  stayOpen = false,
  children,
  ...baseProps
}) => {
  const uniqueId = 'accordion-' + 0; /// TODO
  const accordionClasses = classNames('accordion', {
    'accordion-flush': isFlush,
  });
  return (
    <div className={accordionClasses} {...baseProps} id={uniqueId}>
      {children}
    </div>
  );
};

interface AccordionItemProps extends React.ComponentPropsWithoutRef<'div'> {
  title: string;
  body: string | HTMLElement;
  isCollapsed?: boolean;
}

const AccordionItem: React.FC<AccordionItemProps> = ({
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

export { Accordion, AccordionItem };

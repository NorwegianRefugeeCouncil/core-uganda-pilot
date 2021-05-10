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
  parentId: string;
  uniqueKey: string;
  title: string;
  body: string | HTMLElement;
  isCollapsed?: boolean;
}

const AccordionItem: React.FC<AccordionItemProps> = ({
  parentId,
  uniqueKey,
  title,
  body,
  isCollapsed = true,
  ...baseProps
}) => {
  const headerId = 'header-' + uniqueKey;
  const collapseId = 'collapse-' + uniqueKey;
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
          data-bs-toggle="collapse"
          data-bs-target={'#' + collapseId}
          aria-expanded={isCollapsed}
          aria-controls={collapseId}
        >
          {title}
        </button>
      </h2>
      <div
        id={collapseId}
        className={collapseClass}
        aria-labelledby={headerId}
        data-bs-parent={'#' + parentId}
      >
        <div className="accordion-body">{body}</div>
      </div>
    </>
  );
};

export { Accordion, AccordionItem };

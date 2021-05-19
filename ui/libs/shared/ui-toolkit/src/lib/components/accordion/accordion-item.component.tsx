import * as React from 'react';
import classNames from 'classnames';

export interface AccordionItemProps
  extends React.ComponentPropsWithoutRef<'div'> {
  title: string;
  body: string | HTMLElement;
  isCollapsed?: boolean;
}

export const AccordionItem: React.FC<AccordionItemProps> = ({
  title,
  body,
  isCollapsed = true,
  ...baseProps
}) => {
  const buttonClass = classNames('accordion-button', {
    collapsed: isCollapsed,
  });
  const collapseClass = classNames('accordion-collapse collapse', {
    show: !isCollapsed,
  });
  return (
    <>
      <h2 className="accordion-header">
        <button
          className={buttonClass}
          type="button"
          aria-expanded={isCollapsed}
        >
          {title}
        </button>
      </h2>
      <div className={collapseClass}>
        <div className="accordion-body">{body}</div>
      </div>
    </>
  );
};

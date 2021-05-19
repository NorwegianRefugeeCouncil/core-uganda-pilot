import * as React from 'react';

export const AccordionBody = (props) => {
  return (
    <div className="accordion-body" {...props}>
      {props.children}
    </div>
  );
};

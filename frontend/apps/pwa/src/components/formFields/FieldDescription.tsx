import React from 'react';

type Props = {
  text: string;
  fieldId: string;
};

export const FieldDescription: React.FC<Props> = ({ text, fieldId }) => {
  if (text) {
    return (
      <small className="text-muted" id={`description-${fieldId}`}>
        {text}
      </small>
    );
  }
  return <></>;
};

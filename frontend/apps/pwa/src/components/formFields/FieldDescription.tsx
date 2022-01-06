import React from 'react';

type Props = {
  text: string;
};

export const FieldDescription: React.FC<Props> = ({ text }) => {
  if (text) {
    return <small className="text-muted">{text}</small>;
  }
  return <></>;
};

import * as React from 'react';
import './formrenderer.module.css';
import { getSchema } from './get-schema';
import { parseSchema } from './parse-schema';

export const FormRenderer = (props) => {
  const [schema, setSchema] = React.useState(null);

  React.useEffect(() => {
    getSchema().then((data) => setSchema(data));
  });

  return <form {...props}>{parseSchema(schema)}</form>;
};

import * as React from 'react';
import { FormContext, FormContextInterface } from './form-context';

export interface FormGroupProps extends React.ComponentPropsWithoutRef<'div'> {
  controlId: string;
}

export const FormGroup: React.FC<FormGroupProps> = ({ controlId, ...rest }) => {
  const formCtx: FormContextInterface = { controlId };
  const className = 'mb-3';
  return (
    <FormContext.Provider value={formCtx}>
      <div className={className} {...rest} />
    </FormContext.Provider>
  );
};

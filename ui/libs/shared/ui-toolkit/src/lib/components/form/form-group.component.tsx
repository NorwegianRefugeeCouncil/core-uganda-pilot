import * as React from 'react';
import FormContext, { FormContextInterface } from './form-context';

interface Props<C extends React.ElementType> {
  as?: C;
  children: React.ReactNode;
  controlId?: string;
}

type FormGroupProps<C extends React.ElementType> = Props<C> &
  Omit<React.ComponentPropsWithoutRef<C>, keyof Props<C>>;

const FormGroup = <C extends React.ElementType = 'div'>({
  as,
  controlId,
  children,
  ...rest
}) => {
  const formCtx: FormContextInterface = { controlId };
  const Component = as ?? 'div';
  return (
    <FormContext.Provider value={formCtx}>
      <Component {...rest}></Component>
    </FormContext.Provider>
  );
};

export default FormGroup;

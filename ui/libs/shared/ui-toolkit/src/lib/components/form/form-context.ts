import { createContext } from 'react';

export interface FormContextInterface {
  controlId: string;
}

const FormContext = createContext<FormContextInterface | null>(null);

export default FormContext;

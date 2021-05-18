import { createContext } from 'react';

export interface FormContextInterface {
  controlId: string;
}

export const FormContext = createContext<FormContextInterface | null>(null);

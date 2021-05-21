import * as React from 'react';
import { Form } from '@core/ui-toolkit';
import {
  FormDataField,
  FormElementOptions,
  FormElement,
  FormElementType,
  FormSchema,
} from './form-types';

const locale = 'en'; // TODO

interface FormElementProps {
  type: FormElementType;
  options: FormElementOptions;
}

function FormComponent({ type, options }: Omit<FormElement, 'children'>) {
  const name = options.name[locale];
  const description = options.name[locale];
  const tooltip = options.tooltip[locale];
  switch (type) {
    case 'shortText':
      return (
        <Form.Group controlId="shortText">
          <Form.Label>{name}</Form.Label>
          <Form.Control type="text" placeholder={description} />
          <Form.Text>{tooltip}</Form.Text>
        </Form.Group>
      );
    case 'longText':
      return (
        <Form.Group controlId="shortText">
          <Form.Label>{name}</Form.Label>
          <Form.Control type="textarea" placeholder={description} />
          <Form.Text>{tooltip}</Form.Text>
        </Form.Group>
      );
    case 'checkbox':
      return <Form.Check id="checkbox" label={name} />;
    case 'section':
    case 'select':
    case 'date':
    case 'dateTime':
    case 'integer':
    case 'time':
    default:
      return null;
  }
}

function parseFormElement({ type, children, options }: FormElement) {
  if (children == null || children.length === 0) {
    return <FormComponent type={type} options={options} />;
  } else if (children.length) {
    return children.map(parseFormElement);
  }
}

function parseFormData(data: FormDataField) {
  const elements: JSX.Element[] = [];
  for (const field in data) {
    const element = parseFormElement(data[field]);
    elements.push(element);
  }
  return elements;
}

export function parseSchema(schema: FormSchema): JSX.Element[] {
  if (schema == null) return null;

  const form = parseFormData(schema.data);

  return form;
}

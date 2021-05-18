import { storiesOf } from '@storybook/react';
import React from 'react';
import { withKnobs, text } from '@storybook/addon-knobs';
import { Field, FieldOptions, FieldType } from './field.component';

storiesOf('Field Builder UI', module)
  .addDecorator(withKnobs)
  .add('Demo - String', () => (
    <>
      <Field
        key="demo-string"
        type={FieldType.string}
        children={[]}
        options={
          {
            name: {
              en: 'demo string field',
            },
            description: {
              en: 'demo string field description',
            },
            tooltip: {
              en: 'this is a demo string field',
            },
            value: '',
          } as FieldOptions
        }
      />
    </>
  ))
  .add('Demo - Integer', () => (
    <>
      <Field
        key="demo-integer"
        type={FieldType.integer}
        children={[]}
        options={
          {
            name: {
              en: 'demo integer field',
              es: 'campo entero demo',
              fr: 'champ entier de démo',
              ar: 'العدد الصحيح التجريبي'
            },
            description: {
              en: 'demo integer field description',
              ar: 'وصف حقل العدد الصحيح'
            },
            tooltip: {
              en: 'this is a demo integer field',
            },
            value: 0,
          } as FieldOptions
        }
      />
    </>
  ))
  .add('Demo - Float', () => (
    <>
      <Field
        key="demo-float"
        type={FieldType.float}
        children={[]}
        options={
          {
            name: {
              en: 'demo float field',
            },
            description: {
              en: 'demo float field description',
            },
            tooltip: {
              en: 'this is a demo float field',
            },
            value: 1.1,
          } as FieldOptions
        }
      />
    </>
  ))
  .add('Demo - Checkbox', () => (
    <>
      <Field
        key="demo-checkbox"
        type={FieldType.checkbox}
        children={[]}
        options={
          {
            name: {
              en: 'demo checkbox field',
            },
            description: {
              en: 'demo checkbox field description',
            },
            tooltip: {
              en: 'this is a demo checkbox field',
            },
            value: false,
          } as FieldOptions
        }
      />
    </>
  ))
  .add('Demo - Radio', () => (
    <>
      <Field
        key="demo-radio"
        type={FieldType.radio}
        children={[]}
        options={
          {
            name: {
              en: 'demo radio field',
            },
            description: {
              en: 'demo radio field description',
            },
            tooltip: {
              en: 'this is a demo radio field',
            },
            value: false,
          } as FieldOptions
        }
      />
    </>
  ))
  .add('Demo - Select', () => (
    <>
      <Field
        key="demo-select"
        type={FieldType.select}
        children={[]}
        options={
          {
            name: {
              en: 'demo select field',
            },
            description: {
              en: 'demo select field description',
            },
            tooltip: {
              en: 'this is a demo select field',
            },
            value: "",
          } as FieldOptions
        }
      />
    </>
  ))
  .add('Demo - Multiselect', () => (
    <>
      <Field
        key="demo-multiselect"
        type={FieldType.multiselect}
        children={[]}
        options={
          {
            name: {
              en: 'demo multiselect field',
            },
            description: {
              en: 'demo multiselect field description',
            },
            tooltip: {
              en: 'this is a demo multiselect field',
            },
            value: [],
          } as FieldOptions
        }
      />
    </>
  ))

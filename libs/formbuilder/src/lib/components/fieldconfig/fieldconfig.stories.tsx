import { storiesOf } from '@storybook/react';
import React from 'react';
import { FieldConfig, GenericFieldConfig } from './fieldconfig.component'
import { withKnobs, text } from '@storybook/addon-knobs';
import { FieldType } from '../fieldtype/fieldtype.component';

storiesOf('Field Config', module)
.addDecorator(withKnobs)
.add('default', () => (
  <>
    <FieldConfig fieldType={text('Type', FieldType.text) as FieldType} fieldProps={{} as GenericFieldConfig}/>
  </>
));
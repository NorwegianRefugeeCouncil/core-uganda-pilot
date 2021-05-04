import { storiesOf } from '@storybook/react';
import React from 'react';
import { FieldTypePicker } from './fieldtype.component'
import { withKnobs, text } from '@storybook/addon-knobs';

storiesOf('Field Type Picker', module)
.addDecorator(withKnobs)
.add('default', () => (
  <>
    <FieldTypePicker value={undefined} />
  </>
));
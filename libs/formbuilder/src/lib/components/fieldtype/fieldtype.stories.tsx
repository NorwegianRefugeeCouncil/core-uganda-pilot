import { storiesOf } from '@storybook/react';
import React from 'react';
import { FieldTypePicker } from './fieldtype.component'
import {
    Card
} from '@nrc.no/ui-toolkit'

storiesOf('Field Type Picker', module).add('default', () => (
  <>
    <FieldTypePicker value={undefined} />
  </>
));
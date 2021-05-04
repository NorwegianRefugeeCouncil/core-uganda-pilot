import { storiesOf } from '@storybook/react';
import React from 'react';
import { FieldConfig } from './fieldconfig.component'
import {
    Card
} from '@nrc.no/ui-toolkit'

storiesOf('Field Config', module).add('default', () => (
  <>
    <FieldConfig />
  </>
));
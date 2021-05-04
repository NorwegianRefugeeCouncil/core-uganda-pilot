import { storiesOf } from '@storybook/react';
import React from 'react';
import { FieldBuilder } from './fieldbuilder.component'
import {
    Card
} from '@nrc.no/ui-toolkit'

storiesOf('Field Builder', module).add('default', () => (
  <>
    <Card>
      <FieldBuilder />
    </Card>
  </>
));
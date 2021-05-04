import { storiesOf } from '@storybook/react';
import React from 'react';
import { Test } from './test'
import {
    Card
} from '@nrc.no/ui-toolkit'

storiesOf('Test', module).add('default', () => (
  <>
    <Card>
      <Test />
    </Card>
  </>
));